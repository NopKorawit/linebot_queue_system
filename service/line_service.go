package service

import (
	"errors"
	"fmt"
	"line/model"
	"line/repository"
	"log"
)

type queueService struct {
	queueRepo repository.QueueRepository //อ้างถึง interface
}

//constructor
func NewQueueService(queueRepo repository.QueueRepository) QueueService {
	return queueService{queueRepo: queueRepo}
}

func (s queueService) GetQueue(code string) (*model.QueueResponseLine, error) {
	queue, err := s.queueRepo.GetQueuesByCode(code)
	if err != nil {
		log.Println(err)
		return nil, errors.New("repository error")
	}
	current, err := s.queueRepo.GetCurrentQueue(queue.Type)
	if err != nil {
		log.Println(err)
		return nil, errors.New("repository error")
	}
	qReponse := model.QueueResponseLine{
		CurrentCode: fmt.Sprintf("%v%03d", current.Type, current.Code),
		UserCode:    fmt.Sprintf("%v%03d", queue.Type, queue.Code),
		QueueAmount: queue.Code - current.Code,
		Date:        queue.Date,
		Name:        queue.Name,
	}
	fmt.Println(qReponse)
	return &qReponse, nil
}

func (s queueService) DeleteQueuebyUID(UserID string) error {
	queue, err := s.queueRepo.DeleteQueuebyUID(UserID)
	if err != nil {
		log.Println(err)
		return errors.New("repository error")
	}
	log.Printf("%v is cancle queue", queue.Name)
	return nil
}

// func pushmessage (userID string,message string){
// 	bot, err := linebot.New(<channel secret>, <channel token>)
// 	if err != nil {
// 	...
// 	}
// 	if _, err := bot.Multicast(userIDs, linebot.NewTextMessage("hello")).Do(); err != nil {
// 	...
// 	}
// }
