package Egoconf

import "time"

const ENV = "dev"

//const ENV = "test"

//const ENV = "prod"

var ConfMap = make(map[string]Conf)

type Conf struct {
	LogPath    string
	LogName    string
	VerifyRole SmsRole
}

type SmsRole struct {
	/**
	 * 阿里短信通道开关
	 */
	AliChannel bool

	/**
	 * 云之讯短信通道开关
	 */
	YzxChannel bool
	/**
	 * 同一账号每天全业务短信验证码总量(ali) 控制台可更改
	 */
	AliSendLimitCount int

	/**
	 * 同一账号每天全业务短信验证码总量(yzx)
	 */
	YzxSendLimitCount int
}

//dev 环境
var DevConf = Conf{
	//必须以"/"结尾
	LogPath:    "/tmp/sms-service/",
	LogName:    time.Now().Format("20060102") + ".log",
	VerifyRole: SmsRole{AliChannel: true, YzxChannel: true, AliSendLimitCount: 10, YzxSendLimitCount: 8},
}

//test 环境
var TestConf = Conf{
	//必须以"/"结尾
	LogPath:    "/tmp/sms-service/",
	LogName:    time.Now().Format("20060102") + ".log",
	VerifyRole: SmsRole{AliChannel: true, YzxChannel: true, AliSendLimitCount: 10, YzxSendLimitCount: 8},
}

//prod 环境
var ProdConf = Conf{
	//必须以"/"结尾
	LogPath:    "/tmp/sms-service/",
	LogName:    time.Now().Format("20060102") + ".log",
	VerifyRole: SmsRole{AliChannel: true, YzxChannel: true, AliSendLimitCount: 10, YzxSendLimitCount: 8},
}

func init() {
	switch ENV {
	case "dev":
		ConfMap[ENV] = DevConf
	case "test":
		ConfMap[ENV] = TestConf
	case "prod":
		ConfMap[ENV] = ProdConf
	default:
		ConfMap[ENV] = DevConf
	}
}

func GetConfig() Conf {
	return ConfMap[ENV]
}
