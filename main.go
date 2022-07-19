package main

import (
	"line/handler"
	"line/repository"
	"line/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	handler.InitAll()
	//connect to database + auto migrate
	db := handler.ConnectDatabase()
	// queueRepo := repository.NewQueueRepositoryMock2()

	queueRepo := repository.NewQueueRepositoryDB(db)
	queueService := service.NewQueueService(queueRepo)
	queueHandler := handler.NewQueueHandler(queueService)

	route := gin.Default()
	route.Use(cors.Default())

	//Routes
	route.GET("/", queueHandler.Hello)
	route.POST("/callback", queueHandler.Callback)

	//Run Server
	route.Run(":5000")
}
