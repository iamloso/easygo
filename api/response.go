package Egoapi

import (
	"encoding/json"
	"net/http"
)

//var Desc = make(map[int]string)
const E_SUCCESS = 1000                        //成功
const E_REF_LACK = 1001                       //参数缺失
const E_OTHER = 1003                          //其他原因
const E_FAIL = 1004                           //失败
const E_BUSINESS_LIMIT_CONTROL = 1005         //触发分钟级限流
const E_INVALID_PARAMETERS = 1006             //参数异常
const E_MOBILE_NUMBER_ILLEGAL = 1007          //非法手机号
const E_TEMPLATE_MISSING_PARAMETERS = 1008    //模板缺少变量
const E_BLACK_KEY_CONTROL_LIMIT = 1009        //黑名单管控
const E_PARAM_LENGTH_LIMIT = 1010             //参数超出长度限制
const E_TEMPLATE_PARAMS_ILLEGAL = 1011        //模板变量中包含非法关键字
const E_IP_WHITE_LIST = 1012                  //ip频率限制
const E_INVALID_CHAR = 1021                   // 包含非中文字符
const E_OVER_LENGTH = 1022                    // 超过四个中文字符
const E_SENSITIVE_WORD = 1023                 // 您输入的内容含有敏感词汇，请修改重试
const E_SIGN_EXIST = 1024                     // 签名已存在
const E_SIGN_CHECK_FAIL = 1025                // 签名校验失败
const E_UA_GET_FAIL = 1026                    // UA获取失败
const E_EXCEED_SEND_LIMIT_COUNT = 2001        //验证码超出发送次数
const E_CHECK_VERIFYCODE_FAILED = 2002        //检验验证码失败
const E_USER_NOT_EXISTS = 2003                //用户不存在
const E_VERIFY_CODE_INVALID_OR_CHECKED = 2004 //验证码校验失效或已验证
const E_VERIFY_CODE_WRONG = 2020              //验证码错误
const E_VERIFY_CODE_TIMEOUT = 2021            //验证码已失效
const E_VERIFY_CODE_INVALID = 2022            //验证时已使用
const E_SYSTEM = 500                          //Internal Server Error

var DescMap = map[int]string{
	E_SUCCESS:                        "成功",
	E_REF_LACK:                       "参数缺失",
	E_OTHER:                          "其他原因",
	E_FAIL:                           "失败",
	E_BUSINESS_LIMIT_CONTROL:         "触发分钟级限流",
	E_INVALID_PARAMETERS:             "参数异常",
	E_MOBILE_NUMBER_ILLEGAL:          "非法手机号",
	E_TEMPLATE_MISSING_PARAMETERS:    "模板缺少变量",
	E_BLACK_KEY_CONTROL_LIMIT:        "黑名单管控",
	E_PARAM_LENGTH_LIMIT:             "参数超出长度限制",
	E_TEMPLATE_PARAMS_ILLEGAL:        "模板变量中包含非法关键字",
	E_IP_WHITE_LIST:                  "ip频率限制",
	E_INVALID_CHAR:                   "包含非中文字符",
	E_OVER_LENGTH:                    "超过四个中文字符",
	E_SENSITIVE_WORD:                 "您输入的内容含有敏感词汇，请修改重试",
	E_SIGN_EXIST:                     "签名已存在",
	E_SIGN_CHECK_FAIL:                "签名校验失败",
	E_UA_GET_FAIL:                    "UA获取失败",
	E_EXCEED_SEND_LIMIT_COUNT:        "验证码超出发送次数",
	E_CHECK_VERIFYCODE_FAILED:        "检验验证码失败",
	E_USER_NOT_EXISTS:                "用户不存在",
	E_VERIFY_CODE_INVALID_OR_CHECKED: "验证码校验失效或已验证",
	E_VERIFY_CODE_WRONG:              "验证码错误",
	E_VERIFY_CODE_TIMEOUT:            "验证码已失效",
	E_VERIFY_CODE_INVALID:            "验证时已使用",
	E_SYSTEM:                         "Internal Server Error",
}

var Data interface{}

type JsonResponse struct {
	// Reserved field to add some meta information to the API response
	Code int         `json:"code"`
	Desc string      `json:"desc"`
	Data interface{} `json:"data"`
}

// Writes the response as a standard JSON response with StatusOK
func Response(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if Data == nil {
		Data = ""
	}
	if err := json.NewEncoder(w).Encode(&JsonResponse{Code: code, Desc: DescMap[code], Data: Data}); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func ErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JsonResponse{Code: errorCode, Desc: errorMsg, Data: ""})
}
