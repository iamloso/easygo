package Egolib

import (
	"crypto/md5"
	"easygo/conf"
	"encoding/hex"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/xid"
	"regexp"
	"strings"
	"time"
)

var AppKEY = map[string]string{"1000001": "limishuxue.com", "1000002": "limilaoshi.com", "1000003": "limixuexi.com", "1000004": "m.limishuxue.com", "10000000": "limi"}

var RedisConf = Egoconf.GetRedisConf()

/**
 * 签名校验
 */
func SignVerify(sign string, buildSign string) bool {
	InitRedisPool(RedisConf)
	defer RedisPool.Close()

	rc := RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	signArr, err := redis.Strings(rc.Do("SMEMBERS", "verifycode-sign:sign"))
	if err != nil {
		return false
	}
	var signFlag = false
	for _, val := range signArr {
		if val == sign {
			signFlag = true
		}

	}
	if signFlag == true {
		return false
	} else {
		_, err := rc.Do("SADD", "verifycode-sign:sign", sign)
		if err != nil {
			return false
		}
		t := time.Now()
		expireTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix() - t.Unix()
		rc.Do("expire", "verifycode-sign:sign", expireTime)
	}
	return true
}

/**
 * 生成签名
 */
func CreateSign(appid string, phone string, userType string) string {
	appkey := AppKEY[appid]
	x := md5.New()
	x.Write([]byte(appid + appkey + phone + userType))
	return hex.EncodeToString(x.Sum(nil))
}

/**
 * 生成随机id
 */
func GenXid() string {
	id := xid.New()
	return id.String()
}

/**
 * 根据user agent 判断设备名称
 */
func GetDevice(agent string) string {
	var device string
	if match, err := regexp.MatchString(`.+Windows.+`, agent); match == true && err == nil {
		device = "web设备"
	} else if match, err := regexp.MatchString(`.+Macintosh.+`, agent); match == true && err == nil {
		device = "mac设备"
	} else if match, err := regexp.MatchString(`.+iPad.+`, agent); match == true && err == nil {
		device = "iPad设备"
	} else if match, err := regexp.MatchString(`.+iPhone.+`, agent); match == true && err == nil {
		device = "iPhone设备"
	} else if match, err := regexp.MatchString(`.+Android.+`, agent); match == true && err == nil {
		device = "Android设备"
	}

	return device
}

/**
 * ip 请求限制
 * 校验一个ip， 24小时内， 最高请求20次
 */
func IpVerify(remoteIp string) bool {
	innerIp := map[string]bool{"125.33.127.191": true, "125.34.2.84": true, "123.116.42.38": true, "123.115.64.72": true}
	if remoteIp == "" {
		return false
	}
	ip := remoteIp[0:strings.Index(remoteIp, ":")]

	InitRedisPool(RedisConf)
	defer RedisPool.Close()

	rc := RedisPool.Get()
	defer rc.Close()

	rc.Do("SELECT", RedisConf.SelectDb)

	ipNum, _ := redis.Int(rc.Do("GET", "verifycode-ip:"+ip))
	if innerIp[ip] != true && ipNum >= 50 {
		return false
	}
	if ipNum == 0 {
		t := time.Now()
		expireTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix() - t.Unix()
		rc.Do("SETEX", "verifycode-ip:"+ip, expireTime, 1)
	} else {
		rc.Do("INCR", "verifycode-ip:"+ip)
	}
	return true
}
