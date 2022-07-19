package repository

import (
	"errors"
	"line/model"
	"log"
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

// func (r queueRepositoryDB) GetQueuesByCode2(strcode string) (*model.QueueModel, error) {
// 	queue := model.QueueModel{}
// 	num := strings.TrimLeft(strcode, "ABCD")
// 	code, _ := strconv.Atoi(num)
// 	Type := strings.Trim(strcode, num)
// 	err := r.db.Where("Code = ? AND Type = ?", code, Type).First(&queue).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &queue, nil
// }

// func (r queueRepositoryDB) GetCurrentQueue2(types string) (*model.QueueModel, error) {
// 	currentqueue := model.QueueModel{}
// 	err := r.db.Order("Date").Where("Type = ?", types).First(&currentqueue).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &currentqueue, nil
// }

func (r queueRepositoryDB) GetQueuesByCode(strcode string) (*model.QueueModel, error) {
	queue := model.QueueModel{}
	num := strings.TrimLeft(strcode, "ABCD")
	code, _ := strconv.Atoi(num)
	Type := strings.Trim(strcode, num)
	result := r.db.Where("Code = ? AND Type = ?", code, Type).Find(&queue)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user Code not found")
	}
	return &queue, nil
}

func (r queueRepositoryDB) GetCurrentQueue(types string) (*model.QueueModel, error) {
	currentqueue := model.QueueModel{}
	result := r.db.Order("Date").Where("Type = ?", types).Find(&currentqueue)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("current Code not found")
	}
	return &currentqueue, nil
}

func (r queueRepositoryDB) DeleteQueuebyUID(UserID string) (*model.QueueModel, error) {
	queue := model.QueueModel{}
	result := r.db.Order("Date").Where("user_id = ?", UserID).Find(&queue)
	log.Println(&queue)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("current Code not found")
	}
	r.db.Where("user_id", UserID).Delete(&queue)
	return &queue, nil
}
