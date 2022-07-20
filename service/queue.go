package service

import (
	"line/model"
)

//port
type QueueService interface {
	GetQueue(Code string) (*model.QueueResponseLine, error)
	DeleteQueuebyUID(UserID string) error
}
