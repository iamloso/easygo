package Egomodels

import (
	"easygo/conf"
	"easygo/lib"
	"easygo/log"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	//"ms-go-sdk/aliyun-ms"
	//"ms-go-sdk/yzx-ms"
	"strconv"
	"time"
)

const (
	CHANNEL_ALI            = "ali"
	CHANNEL_YZX            = "yzx"
	CHANNEL_ALI_VOICE      = "voice"
	ALI_TEMPLATE_CODE      = "120410015"
	YZX_TEMPLATE_CODE      = "14792"
	ALIVOICE_TEMPLATE_CODE = "116565438"
)

var minuteKey = "sms:phone:businesstype:minute:"
var dayKey = "sms:phone:businesstype:day:"
var verifyCodeKey = "sms:verify:code:"
var VCode string

var VerifyRole Egoconf.SmsRole

var AliChannel aliyunms.Sms
var AliVoice aliyunms.Vms
var YzxChannel yzxms.Sms

type VerifyCode struct {
	Id            int       `orm:"column(id);auto"`
	Phone         string    `orm:"column(phone);size(11)"`
	VerifyCode    string    `orm:"column(verify_code);size(4)"`
	TempleteId    int       `orm:"column(templete_id);null"`
	Product       string    `orm:"column(product);size(10)"`
	IsValid       int8      `orm:"column(is_valid);null" description:"短信状态：-1未发送0待发送1已发送2已验证3已无效（已使用）"`
	AddTime       time.Time `orm:"column(add_time);type(timestamp);auto_now_add"`
	Type          int       `orm:"column(type)" description:"验证码类型值：0:未分组 1学生注册 2老师注册 3学生忘记密码 4老师忘记密码 5学生验证码登录 6老师验证码登录 7狸米课堂 8狸米家长注册"`
	ChannelType   int8      `orm:"column(channel_type)" description:"渠道类型，1普通短信 2语言验证码"`
	Losetime      int       `orm:"column(losetime)" description:"有效期（单位 秒）"`
	BusinessType  int8      `orm:"column(business_type);null" description:"业务类型：0默认值，1学习"`
	PhoneOperator string    `orm:"column(phone_operator);size(55)" description:"手机运营商"`
	PhoneCity     string    `orm:"column(phone_city);size(55)" description:"手机归属地区"`
	BizId         string    `orm:"column(biz_id);size(30)" description:"发送流水号"`
	RequestId     string    `orm:"column(request_id);size(50)" description:"发送短信请求id"`
	SentMsg       string    `orm:"column(sent_msg);size(55)" description:"平台返回的状态信息"`
	SentCode      string    `orm:"column(sent_code);size(55)" description:"平台返回的状态码"`
	ReportTime    time.Time `orm:"column(report_time);type(timestamp);auto_now_add" description:"终端接收时间"`
	SmsSize       int8      `orm:"column(sms_size)" description:"短信长度"`
	CheckNum      int8      `orm:"column(check_num)" description:"校验次数"`
	ReportMsg     string    `orm:"column(report_msg);size(80)" description:"终端状态错误信息"`
	ReportCode    string    `orm:"column(report_code);size(55)" description:"终端状态错误码"`
}

var RedisConf = Egoconf.GetRedisConf()

func init() {
	VerifyRole = Egoconf.GetConfig().VerifyRole
	orm.RegisterModel(new(VerifyCode))

	Egolog.Info("verifyCode model 初始化完成!")
}

func (t *VerifyCode) TableName() string {
	return "verify_code"
}

// AddVerifyCode insert a new VerifyCode into database and returns
// last inserted Id on success.
func AddVerifyCode(m *VerifyCode) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetVerifyCodeById retrieves VerifyCode by Id. Returns error if
// Id doesn't exist
func GetVerifyCodeById(id int) (v *VerifyCode, err error) {
	o := orm.NewOrm()
	v = &VerifyCode{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateVerifyCode updates VerifyCode by Id and returns error if
// the record to be updated doesn't exist
func UpdateVerifyCodeById(m *VerifyCode) (err error) {
	o := orm.NewOrm()
	v := VerifyCode{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

/**
 * 短信数据存储redis
 */
func BuildRedisData(phone string, businessType string, channelRoute string) bool {
	Egolog.InfoData("BuildRedisData：(开启跟踪)短信数据存储redis", map[string]interface{}{"phone": phone, "businessType": businessType})

	SetMinuteLevel(phone, businessType)

	SetDayLevel(phone, businessType, channelRoute)

	return true
}

/**
 * 分钟级别， redis 数据
 */
func SetMinuteLevel(phone string, businessType string) bool {
	Egolog.InfoData("SetMinuteLevel：(开启跟踪)分钟级别短信数据存储redis", map[string]interface{}{"phone": phone, "businessType": businessType})

	Egolib.InitRedisPool(RedisConf)
	defer Egolib.RedisPool.Close()

	var rc = Egolib.RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	if _, err := rc.Do("SETEX", minuteKey+phone+":"+businessType, 60, 1); err != nil {
		Egolog.ErrorData("SetMinuteLevel：(系统错误)分钟级别短信数据存储失败！", map[string]interface{}{"phone": phone, "businessType": businessType})
		return false
	}

	if _, err := rc.Do("SETEX", verifyCodeKey+phone+":"+businessType, 1800, VCode); err != nil {
		Egolog.ErrorData("SetMinuteLevel：(系统错误)分钟级别短信数据存储失败！", map[string]interface{}{"phone": phone, "businessType": businessType})
		return false
	}

	return true
}

/**
 * 分钟级别数据是否有效
 */
func GetMinuteLevel(phone string, businessType string) bool {
	Egolib.InitRedisPool(RedisConf)
	defer Egolib.RedisPool.Close()

	var rc = Egolib.RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	flag, _ := redis.Bool(rc.Do("GET", minuteKey+phone+":"+businessType))
	return flag
}

/**
 * 天级别， redis 数据
 */
func SetDayLevel(phone string, businessType string, channelRoute string) bool {
	Egolog.InfoData("SetDayLevel：(开启跟踪)天级别短信统计数据存储redis", map[string]interface{}{"phone": phone, "businessType": businessType})

	Egolib.InitRedisPool(RedisConf)
	defer Egolib.RedisPool.Close()

	var rc = Egolib.RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	t := time.Now()
	expireTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix() - t.Unix()

	aliSendNum, yzxSendNum := GetDayLevel(phone, businessType)

	var err error
	if channelRoute == CHANNEL_ALI {
		if aliSendNum > 0 {
			_, err = rc.Do("INCR", dayKey+"aliyun:"+phone+":"+businessType)
		} else {
			_, err = rc.Do("SETEX", dayKey+"aliyun:"+phone+":"+businessType, expireTime, 1)
		}
		if err != nil {
			Egolog.ErrorData("SetDayLevel：(系统错误)天级别短信统计数据存储失败(aliyun)", map[string]interface{}{"phone": phone, "businessType": businessType})
			return false
		}
	}
	if channelRoute == CHANNEL_YZX {
		if yzxSendNum > 0 {
			_, err = rc.Do("INCR", dayKey+"yzx:"+phone+":"+businessType)
		} else {
			_, err = rc.Do("SETEX", dayKey+"yzx:"+phone+":"+businessType, expireTime, 1)
		}
		if err != nil {
			Egolog.ErrorData("SetDayLevel：(系统错误)天级别短信统计数据存储失败(yzx)", map[string]interface{}{"phone": phone, "businessType": businessType})
			return false
		}
	}

	if channelRoute == CHANNEL_ALI_VOICE {
		if yzxSendNum > 0 {
			_, err = rc.Do("INCR", dayKey+"aliyun_voice:"+phone+":"+businessType)
		} else {
			_, err = rc.Do("SETEX", dayKey+"aliyun_voice:"+phone+":"+businessType, expireTime, 1)
		}
		if err != nil {
			Egolog.ErrorData("SetDayLevel：(系统错误)天级别短信统计数据存储失败(aliyun_voice)", map[string]interface{}{"phone": phone, "businessType": businessType})
			return false
		}
	}

	return true
}

/**
 * 获取天级别数据
 */
func GetDayLevel(phone string, businessType string) (int, int) {
	Egolib.InitRedisPool(RedisConf)
	defer Egolib.RedisPool.Close()

	var rc = Egolib.RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	aliSendNum, _ := redis.Int(rc.Do("GET", dayKey+"aliyun:"+phone+":"+businessType))

	yzxSendNum, _ := redis.Int(rc.Do("GET", dayKey+"yzx:"+phone+":"+businessType))

	return aliSendNum, yzxSendNum
}

/**
 * 通道路由
 * 根据当天通道发送量与通道开关选择合适通道
 * 1 通道发送量两者达到最大限制， 选择阿里语音通道
 * 2 一方通道达到最大限制， 如通道开关打开，选择另一方通道， 否则选择阿里语音通道
 */
func ChannelRoute(phone string, businessType string) string {
	Egolog.InfoData("ChannelRoute:(开启跟踪)短信通道路由", map[string]interface{}{"phone": phone, "businessType": businessType})

	aliSendNum, yzxSendNum := GetDayLevel(phone, businessType)

	if aliSendNum >= VerifyRole.AliSendLimitCount && yzxSendNum >= VerifyRole.YzxSendLimitCount {
		return CHANNEL_ALI_VOICE
	} else if aliSendNum >= VerifyRole.AliSendLimitCount && yzxSendNum < VerifyRole.YzxSendLimitCount {
		if VerifyRole.YzxChannel {
			return CHANNEL_YZX
		}
	} else if aliSendNum < VerifyRole.AliSendLimitCount && yzxSendNum >= VerifyRole.YzxSendLimitCount {
		if VerifyRole.AliChannel {
			return CHANNEL_ALI
		}
	} else if aliSendNum < VerifyRole.AliSendLimitCount && yzxSendNum < VerifyRole.YzxSendLimitCount {
		if VerifyRole.AliChannel && VerifyRole.YzxChannel {
			return ChannelRandom()
		}
	}
	return CHANNEL_ALI_VOICE
}

/**
 * 随机短信通道
 */
func ChannelRandom() string {
	weightAli := VerifyRole.AliSendLimitCount * 100 / (VerifyRole.AliSendLimitCount + VerifyRole.YzxSendLimitCount)
	randNum := rand.Intn(100)
	weightAli = 0
	if weightAli >= randNum {
		return CHANNEL_ALI
	}

	return CHANNEL_YZX
}

/**
 * 调起发送短信sdk， 发送验证码
 */
func SendVerifyCode(phone string, businessType string) (string, bool) {
	var params = map[string]interface{}{"phone": phone, "businessType": businessType}
	Egolog.InfoData("SendVerifyCode:(开启跟踪)发送短信验证码", params)

	channel := ChannelRoute(phone, businessType)

	params["channelRoute"] = channel

	Egolog.InfoData("SendVerifyCode:(开启跟踪)通道路由信息", params)

	VCode = strconv.Itoa(GetVerifyCode(phone, businessType))
	var err error
	var AliResponse aliyunms.Response
	var YzxResponse yzxms.Response
	codeData := VerifyCode{}
	codeData.Phone = phone
	codeData.VerifyCode = VCode
	codeData.IsValid = 1
	busiType, _ := strconv.Atoi(businessType)
	codeData.BusinessType = int8(busiType)
	codeData.ChannelType = 1
	codeData.Losetime = 1800
	if channel == CHANNEL_ALI {
		AliResponse, err = AliChannel.Send(phone, "{\"code\":\""+VCode+"\"}", "SMS_"+ALI_TEMPLATE_CODE)
		codeData.Product = CHANNEL_ALI
		codeData.RequestId = AliResponse.RequestId
		codeData.BizId = AliResponse.BizId
		codeData.SentCode = AliResponse.Code
		codeData.SentMsg = AliResponse.Message
		codeData.TempleteId, _ = strconv.Atoi(ALI_TEMPLATE_CODE)

	} else if channel == CHANNEL_YZX {
		YzxResponse, err = YzxChannel.Send(phone, VCode, YZX_TEMPLATE_CODE)
		codeData.Product = CHANNEL_YZX
		codeData.RequestId = YzxResponse.RequestId
		codeData.BizId = YzxResponse.BizId
		codeData.SentCode = YzxResponse.Code
		codeData.SentMsg = YzxResponse.Message
		codeData.TempleteId, _ = strconv.Atoi(YZX_TEMPLATE_CODE)
	} else if channel == CHANNEL_ALI_VOICE {
		AliResponse, err = AliVoice.Send(phone, "{\"code\":\""+VCode+"\"}", "TTS_"+ALIVOICE_TEMPLATE_CODE)
		codeData.Product = CHANNEL_ALI_VOICE
		codeData.RequestId = AliResponse.RequestId
		codeData.BizId = AliResponse.BizId
		codeData.SentCode = AliResponse.Code
		codeData.SentMsg = AliResponse.Message
		codeData.TempleteId, _ = strconv.Atoi(ALIVOICE_TEMPLATE_CODE)
		codeData.ChannelType = 2
	}
	if _, dbErr := AddVerifyCode(&codeData); dbErr != nil {
		params["system_error"] = dbErr
		Egolog.ErrorData("SendVerifyCode:(系统错误)短信数据入库失败！", params)
	}

	if err != nil {
		params["system_error"] = err
		Egolog.ErrorData("SendVerifyCode:(系统错误)短信sdk报错！", params)
		return channel, false
	}
	BuildRedisData(phone, businessType, channel)

	return channel, true
}

/**
 * 生成短信验证码
 */
func makeVerifyCode() int {
	VCode := rand.Intn(9999)
	return VCode
}

/**
 * 获取短信验证码
 */
func GetVerifyCode(phone string, businessType string) int {
	Egolib.InitRedisPool(RedisConf)
	defer Egolib.RedisPool.Close()

	var rc = Egolib.RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	VCode, _ := redis.Int(rc.Do("GET", verifyCodeKey+phone+":"+businessType))
	if VCode == 0 {
		VCode = makeVerifyCode()
	}

	return VCode
}
