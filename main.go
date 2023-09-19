package main

import (
	"github.com/joho/godotenv"
	"log"
	"sukitime.com/v2/bootstrap"
	"sukitime.com/v2/web"
)

func main() {
	var w = make(chan int, 1)

	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file")
	}
	//初始化DB
	bootstrap.InitGorm()
	//初始化Redis
	bootstrap.InitRedis()
	//初始化GRPC
	go bootstrap.InitExportRpcServer()
	//初始化路由
	web.LoadRouter()
	<-w
}
