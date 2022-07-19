package service

import (
	"line/model"
)

//port
type QueueService interface {
	GetQueue(Code string) (*model.QueueResponse, error)
	DeleteQueuebyUID(UserID string) error
}
