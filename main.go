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
					userIDs := "U75d559eb17b924479b63d01491314f48"
					queue, err := queueService.GetQueue(message.Text)
					if err != nil {
						if _, err := bot.PushMessage(userIDs, linebot.NewTextMessage("hello")).Do(); err != nil {
							log.Print(err)
						}
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ไม่พบเลขคิวที่คุณค้นหาหรืออาจเลยคิวของคุณมาแล้ว")).Do(); err != nil {
							log.Print(err)
							return
						}
					}
					var wait string
					if queue.QueueAmount == 1 {
						wait = "Waiting a queue"
					} else if queue.QueueAmount > 1 {
						wait = fmt.Sprintf("Waiting %v queues", queue.QueueAmount)
					} else if queue.QueueAmount == 0 {
						wait = "It's your turn"
					}
					flex := fmt.Sprintf(`{
						"type": "bubble",
						"size": "kilo",
						"direction": "ltr",
						"hero": {
						  "type": "image",
						  "url": "https://pbs.twimg.com/media/DYY92oBVwAAJKvb.jpg",
						  "size": "full",
						  "aspectRatio": "25:13",
						  "aspectMode": "cover",
						  "position": "relative"
						},
						"body": {
						  "type": "box",
						  "layout": "vertical",
						  "spacing": "md",
						  "contents": [
							{
							  "type": "text",
							  "text": "%v",
							  "weight": "bold",
							  "size": "xl",
							  "gravity": "center",
							  "margin": "lg",
							  "wrap": true,
							  "contents": []
							},
							{
							  "type": "box",
							  "layout": "vertical",
							  "spacing": "sm",
							  "margin": "lg",
							  "contents": [
								{
								  "type": "box",
								  "layout": "baseline",
								  "spacing": "sm",
								  "margin": "xs",
								  "contents": [
									{
									  "type": "text",
									  "text": "Name",
									  "size": "sm",
									  "color": "#AAAAAA",
									  "flex": 2,
									  "contents": []
									},
									{
									  "type": "text",
									  "text": "%v",
									  "size": "sm",
									  "color": "#666666",
									  "flex": 4,
									  "wrap": true,
									  "contents": []
									}
								  ]
								},
								{
								  "type": "box",
								  "layout": "baseline",
								  "spacing": "sm",
								  "margin": "xs",
								  "contents": [
									{
									  "type": "text",
									  "text": "Date",
									  "size": "sm",
									  "color": "#AAAAAA",
									  "flex": 2,
									  "contents": []
									},
									{
									  "type": "text",
									  "text": "%v",
									  "size": "sm",
									  "color": "#666666",
									  "flex": 4,
									  "wrap": true,
									  "contents": []
									}
								  ]
								}
							  ]
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "spacing": "sm",
							  "margin": "xs",
							  "contents": [
								{
								  "type": "text",
								  "text": "Queue",
								  "size": "sm",
								  "color": "#AAAAAA",
								  "flex": 2,
								  "contents": []
								},
								{
								  "type": "text",
								  "text": "%v",
								  "size": "sm",
								  "color": "#666666",
								  "flex": 4,
								  "wrap": true,
								  "contents": []
								}
							  ]
							},
							{
							  "type": "box",
							  "layout": "vertical",
							  "margin": "lg",
							  "contents": [
								{
								  "type": "spacer",
								  "size": "xs"
								},
								{
								  "type": "image",
								  "url": "https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=%v",
								  "size": "md",
								  "aspectMode": "cover"
								},
								{
								  "type": "text",
								  "text": "You can enter the restaurant by using this code instead of a ticket",
								  "size": "xxs",
								  "color": "#AAAAAA",
								  "margin": "xxl",
								  "wrap": true,
								  "contents": []
								}
							  ]
							}
						  ]
						}
					  }`, queue.UserCode, queue.Name, queue.Date.Format("Monday 2, 15:04:05"), wait, queue.UserCode)

					// Unmarshal JSON
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flex))
					if err != nil {
						log.Println(err)
					}
					// New Flex Message
					flexMessage := linebot.NewFlexMessage(queue.UserCode, flexContainer)
					// Reply Message
					_, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
					if err != nil {
						log.Print(err)
					}

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
