package db

import (
	"fmt"
	"log"

	"github.com/yoshinori0811/chat_app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {

	url := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`,
		config.Config.DBUserName,
		config.Config.DBUserPassword,
		config.Config.DBHost,
		config.Config.DBPort,
		config.Config.DBName,
	)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connceted")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
