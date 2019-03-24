package Egoconf

type MysqlConf struct {
	Host     string
	UserName string
	PassWord string
	DbName   string
	Port     string
}

var MysqlConfMap = make(map[string]MysqlConf)

var DevMysqlConf = MysqlConf{
	Host:     "127.0.0.1",
	UserName: "root",
	PassWord: "123456",
	DbName:   "test",
	Port:     "3306",
}

var TestMysqlConf = MysqlConf{
	Host:     "127.0.0.1",
	UserName: "root",
	PassWord: "123456",
	DbName:   "test",
	Port:     "3306",
}

var ProdMysqlConf = MysqlConf{
	Host:     "127.0.0.1",
	UserName: "root",
	PassWord: "123456",
	DbName:   "test",
	Port:     "3306",
}

func init() {
	switch ENV {
	case "dev":
		MysqlConfMap[ENV] = DevMysqlConf
	case "test":
		MysqlConfMap[ENV] = TestMysqlConf
	case "prod":
		MysqlConfMap[ENV] = ProdMysqlConf
	default:
		MysqlConfMap[ENV] = DevMysqlConf
	}
}

func GetMysqlConf() MysqlConf {
	return MysqlConfMap[ENV]
}
