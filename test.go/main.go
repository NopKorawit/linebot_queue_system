package main

import (
	"line/handler"
	"line/repository"
	"line/service"
)

func main() {
	handler.InitAll()
	//connect to database + auto migrate
	db := handler.ConnectDatabase()

	//Use Mock Data Repository to test
	// queueRepo := repository.NewQueueRepositoryMock2()

	queueRepo := repository.NewQueueRepositoryDB(db)
	queueService := service.NewQueueService(queueRepo)

	queueService.GetQueue("A004")

}
