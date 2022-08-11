package cmd

import (
	"go-api/api/v1/auth"
	"go-api/api/v1/item"
	"go-api/api/v1/organization"
	"go-api/api/v1/purchaseorder"
	"go-api/api/v1/setting"
	"go-api/api/v1/warehouse"
	"go-api/core/config"
	"go-api/core/database"
	"go-api/core/event"
	"go-api/core/log"
	"go-api/core/router"
)

func Run(args []string) {
	config.LoadConfig(args[1])
	log.ConfigLogger()
	// cache.ConfigCache()
	database.ConfigMysql()
	event.Subscribe(auth.Subscribe)
	r := router.InitRouter()
	router.InitPublicRouter(r, auth.Routers, organization.Routers)
	router.InitAuthRouter(r, auth.AuthRouter, setting.AuthRouter, item.AuthRouter, purchaseorder.AuthRouter, warehouse.AuthRouter)
	router.RunServer(r)
}
