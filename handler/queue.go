package handler

import (
	"fmt"
	"line/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

type QueueHandler interface {
	Callback(c *gin.Context)
	Hello(c *gin.Context)
}

type queueHandler struct {
	qService service.QueueService
}

func NewQueueHandler(qService service.QueueService) QueueHandler {
	return queueHandler{qService: qService}
}

func (h queueHandler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func (h queueHandler) Callback(c *gin.Context) {
	bot := GetBot()
	events, err := bot.ParseRequest(c.Request)
	fmt.Println(err)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				userIDs := "U75d559eb17b924479b63d01491314f48"
				if message.Text == "golf" {
					if _, err := bot.PushMessage(userIDs, linebot.NewTextMessage("มีคนอยากเซ็ทหย่อสูดต่อซูดผ่อซีหม่อสองห่อใส่ไข่กับคุณ")).Do(); err != nil {
						log.Print(err)
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ส่งข้อความให้กอล์ฟแล้ว")).Do(); err != nil {
						log.Print(err)
						return
					}
				}
				queue, err := h.qService.GetQueue(message.Text)
				fmt.Println(err)
				if err != nil {
					if err.Error() == "repository error" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ไม่พบเลขคิวที่คุณค้นหาหรืออาจเลยคิวของคุณมาแล้ว")).Do(); err != nil {
							log.Print(err)
							return
						}
						return
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ระบบผิดพลาด")).Do(); err != nil {
							log.Print(err)
							return
						}
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
				if queue.Name == "" {
					queue.Name = "ไม่ระบุชื่อ"
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
							  "url": "https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=%s",
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
}
