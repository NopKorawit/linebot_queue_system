package handler

import (
	"fmt"
	"line/model"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() (db *gorm.DB) {

	//Set Data source name
	dsn := fmt.Sprintf("server=%v\\%v;Database=%v;praseTime=true",
		viper.GetString("db.server"),
		viper.GetString("db.driver"),
		viper.GetString("db.database"),
	)
	dial := sqlserver.Open(dsn)

	database, err := gorm.Open(dial, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("Failed to connect to database!")
	}
	//auto migration
	database.AutoMigrate(&model.QueueModel{})
	return database
}

func initConfig() {
	//set Read form config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func initTimeZone() {
	//set timezone thailand
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

func InitAll() {
	initTimeZone()
	initConfig()
}
