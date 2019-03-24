package main

import (
	"easygo/conf"
	"easygo/route"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func init() {
	mc := Egoconf.GetMysqlConf()

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", mc.UserName+":"+mc.PassWord+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.DbName+"?charset=utf8&loc=Local")
}

func main() {

	router := Egorouter.NewRouter(Egorouter.RouteList())
	log.Fatal(http.ListenAndServe(":8081", router))
}
