package repository

import (
	"line/model"
)

//Port -Interface
type QueueRepository interface {
	// GetQueuesByType(Type string) ([]model.QueueModel, error)
	// GetQueuesByName(name string, types string) (*model.QueueModel, error)
	GetQueuesByCode(Code string) (*model.QueueModel, error)
	GetCurrentQueue(types string) (*model.QueueModel, error)
}
