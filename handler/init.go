package handler

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() (db *gorm.DB) {

	//Set Data source name
	// dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?&parseTime=True&loc=Local",
	// 	viper.GetString("db.user"),
	// 	viper.GetString("db.pass"),
	// 	viper.GetString("db.host"),
	// 	viper.GetString("db.port"),
	// 	viper.GetString("db.database"),

	// )
	dsnHeroku := os.Getenv("JAWSDB_DSN")
	dsn := fmt.Sprintf("%v?%v", dsnHeroku, "&parseTime=True&loc=Local")
	dial := mysql.Open(dsn)

	database, err := gorm.Open(dial, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("Failed to connect to database!")
	}
	//auto migration
	// database.AutoMigrate(&model.QueueModel{})
	return database
}

func GetBot() (bot *linebot.Client) {
	bot, err := linebot.New(
		// viper.GetString("line.CHANNEL_SECRET"),
		// viper.GetString("line.CHANNEL_TOKEN"),
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

// func ConnectDatabase() (db *gorm.DB) {
// 	//Set Data source name
// 	dsn := fmt.Sprintf("server=%v\\%v;Database=%v;praseTime=true",
// 		viper.GetString("db.server"),
// 		viper.GetString("db.driver"),
// 		viper.GetString("db.database"),
// 	)
// 	dial := sqlserver.Open(dsn)

// 	database, err := gorm.Open(dial, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

// 	if err != nil {
// 		panic("Failed to connect to database!")
// 	}
// 	//auto migration
// 	database.AutoMigrate(&model.QueueModel{})
// 	return database
// }

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
