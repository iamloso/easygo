package Egolib

import (
	"easygo/conf"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisPool *redis.Pool
var redisServer string
var redisAuth string
var idleTimeout = 60 * time.Second

const MAXIDLE = 3
const MAXACTIVE = 1024

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     MAXIDLE,
		MaxActive:   MAXACTIVE,
		IdleTimeout: idleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func InitRedisPool(conf Egoconf.RedisConf) {
	redisServer = conf.Host + ":" + conf.Port
	redisAuth = conf.Auth
	RedisPool = newPool(redisServer, redisAuth)

	return
}

func RedisString(reply interface{}, err1 error) (value string, err2 error) {
	value, err2 = redis.String(reply, err1)
	return
}
