package Egoconf

type RedisConf struct {
	Host     string
	Auth     string
	Port     string
	Lifetime string
	SelectDb string
	Prefix   string
}

var RedisConfMap = make(map[string]RedisConf)

var DevRedisConf = RedisConf{
	Host:     "127.0.0.1",
	Auth:     "",
	Port:     "6379",
	Lifetime: "3600",
	SelectDb: "4",
	Prefix:   "",
}

var TestRedisConf = RedisConf{
	Host:     "127.0.0.1",
	Auth:     "",
	Port:     "6379",
	Lifetime: "3600",
	SelectDb: "4",
	Prefix:   "",
}

var ProdRedisConf = RedisConf{
	Host:     "127.0.0.1",
	Auth:     "",
	Port:     "6379",
	Lifetime: "3600",
	SelectDb: "4",
	Prefix:   "",
}

func init() {
	switch ENV {
	case "dev":
		RedisConfMap[ENV] = DevRedisConf
	case "test":
		RedisConfMap[ENV] = TestRedisConf
	case "prod":
		RedisConfMap[ENV] = ProdRedisConf
	default:
		RedisConfMap[ENV] = DevRedisConf
	}
}

func GetRedisConf() RedisConf {
	return RedisConfMap[ENV]
}
