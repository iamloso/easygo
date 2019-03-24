package Egorouter

import (
	"easygo/api"
	"easygo/log"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func RouteList() Routes {
	routes := Routes{
		Route{"Index", "GET", "/sms/send", Egoapi.Send},
	}
	return routes
}

func NewRouter(routes Routes) *httprouter.Router {

	router := httprouter.New()
	for _, route := range routes {
		var handle httprouter.Handle

		handle = route.HandlerFunc
		handle = Egolog.Logger(handle)

		router.Handle(route.Method, route.Path, handle)
	}

	return router
}
