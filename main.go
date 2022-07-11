package main

import (
	"fmt"
	"line/handler"
	"line/repository"
	"line/service"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v7/linebot/httphandler"
	"github.com/spf13/viper"
)

func main() {
	handler.InitAll()
	//connect to database + auto migrate
	db := handler.ConnectDatabase()
	queueRepo := repository.NewQueueRepositoryDB(db)
	queueService := service.NewQueueService(queueRepo)

	handler, err := httphandler.New(
		viper.GetString("line.CHANNEL_SECRET"),
		viper.GetString("line.CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					queue, err := queueService.GetQueue(message.Text)
					if err != nil {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ไม่พบเลขคิวที่คุณค้นหาหรืออาจเลยคิวของคุณมาแล้ว")).Do(); err != nil {
							log.Print(err)
							return
						}
					}
					reply := fmt.Sprintf("คุณ %v เหลืออีก %v คิว รอสักครู่นะครับ", queue.Name, queue.QueueAmount)
					if queue.QueueAmount == 0 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ถึงคิวของคุณแล้ว รีบมาด่วนเลยครับ")).Do(); err != nil {
							log.Print(err)
							return
						}
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
						log.Print(err)
					}

					// if message.Text == "กรวิชญ์" {
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("คนหล่อเท่")).Do(); err != nil {
					// 		log.Print(err)
					// 	}
					// } else {
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					// 		log.Print(err)
					// 	}
					// }
				}
			}
		}
	})
	http.Handle("/callback", handler)
	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	// if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {

	//ngrok http localhost:5500
	if err := http.ListenAndServe(":5500", nil); err != nil {
		log.Fatal(err)
	}
}
