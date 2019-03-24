package Egoapi

import (
	"easygo/lib"
	"easygo/log"
	//"easygo/models"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

var phone string
var userType string
var businessType string
var verifyCodeRole = "default"
var appid string
var sign string

func init() {
	Egolog.Info("短信验证码接口完成初始化！")
}

/**
 * 发送短信验证码接口
 */
func Send(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	xid := Egolib.GenXid()
	session := r.Header.Get("x-session")
	if session == "" {
		session = Egolib.GenXid()
	}
	Egolog.LogFields["session"] = session
	Egolog.LogFields["trackId"] = xid

	Egolog.Info("发送短信验证码接口:接受发送短信验证码请求")
	Response(w, E_SUCCESS)

	//if code, ok := checkParams(r); ok != nil {
	//	Egolog.InfoData("发送短信验证码接口:参数校验失败", map[string]interface{}{"code": code, "desc": DescMap[code]})
	//	Response(w, code)
	//	return
	//}
	//if ok := Egomodels.GetMinuteLevel(phone, businessType); ok {
	//	Egolog.InfoData("发送短信验证码接口:非法访问，分钟限频", map[string]interface{}{"phone": phone, "businessType": businessType})
	//	Response(w, E_SUCCESS)
	//	return
	//}
	//
	//if channel, ok := Egomodels.SendVerifyCode(phone, businessType); ok {
	//	Egolog.InfoData("发送短信验证码接口:短信发送成功！", map[string]interface{}{"phone": phone, "businessType": businessType})
	//	if channel == Egomodels.CHANNEL_ALI_VOICE {
	//		Data = map[string]bool{"isVoice": true}
	//	} else {
	//		Data = map[string]bool{"isVoice": false}
	//	}
	//	Response(w, E_SUCCESS)
	//} else {
	//	Egolog.InfoData("发送短信验证码接口:短信发送失败！", map[string]interface{}{"phone": phone, "businessType": businessType})
	//	Response(w, E_OTHER)
	//}
	//
	//defer func() {
	//	if rr := recover(); rr != nil {
	//		Egolog.InfoData("发送短信验证码接口:捕获panic错误信息", map[string]interface{}{"phone": phone, "businessType": businessType, "system_panic": rr})
	//	}
	//}()
}

func checkParams(r *http.Request) (code int, err error) {
	r.ParseForm()

	if len(r.Form["phonenum"]) == 0 || len(r.Form["appid"]) == 0 || len(r.Form["sign"]) == 0 || len(r.Form["business_type"]) == 0 {
		return E_REF_LACK, errors.New("false")
	}

	setParams(r)

	buildSign := Egolib.CreateSign(appid, phone, userType)
	if sign[0:32] != buildSign {
		return E_SIGN_CHECK_FAIL, errors.New("false")
	}

	if ok := Egolib.SignVerify(sign, buildSign); !ok {
		return E_SIGN_EXIST, errors.New("false")
	}

	if match, _ := regexp.MatchString(`^1[1-9]\d{9}$`, phone); match != true {
		return E_MOBILE_NUMBER_ILLEGAL, errors.New("false")
	}

	if ok := Egolib.GetDevice(r.Header.Get("user-agent")); ok == "" {
		return E_UA_GET_FAIL, errors.New("false")
	}
	if ok := Egolib.IpVerify(r.RemoteAddr); !ok {
		return E_IP_WHITE_LIST, errors.New("false")
	}

	return E_SUCCESS, nil
}

func setParams(r *http.Request) {
	phone = r.Form["phonenum"][0]
	if len(r.Form["business_type"]) > 0 {
		businessType = r.Form["business_type"][0]
	}
	if len(r.Form["type"]) > 0 {
		userType = r.Form["type"][0]
	}
	appid = r.Form["appid"][0]
	sign = r.Form["sign"][0]
}
