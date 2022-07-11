package repository

import (
	"line/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

//Adapter private
type queueRepositoryDB struct {
	db *gorm.DB
}

//Constructor Public เพื่อ new instance
func NewQueueRepositoryDB(db *gorm.DB) QueueRepository {
	return queueRepositoryDB{db: db}
}

func (r queueRepositoryDB) GetQueuesByCode(strcode string) (*model.QueueModel, error) {
	queue := model.QueueModel{}
	num := strings.TrimLeft(strcode, "ABCD")
	code, _ := strconv.Atoi(num)
	Type := strings.Trim(strcode, num)
	err := r.db.Where("Code = ? AND Type = ?", code, Type).First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r queueRepositoryDB) GetCurrentQueue(types string) (*model.QueueModel, error) {
	currentqueue := model.QueueModel{}
	err := r.db.Order("Date").Where("Type = ?", types).First(&currentqueue).Error
	if err != nil {
		return nil, err
	}
	return &currentqueue, nil
}
