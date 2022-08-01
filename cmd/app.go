package cmd

import (
	"go-api/api/v1/auth"
	"go-api/api/v1/organization"
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
	// router.InitAuthRouter(r, example.Routers)
	router.RunServer(r)
}
