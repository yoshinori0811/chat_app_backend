package db

import (
	"fmt"
	"log"
	"time"

	"github.com/yoshinori0811/chat_app_backend/config"
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

	fmt.Println("url: ", url)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ { // 最大10回再試行
		db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second) // 5秒待ってから再試行
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to get sql.DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = sqlDB.Ping()
		if err == nil {
			log.Println("Connected to DB!")
			return db
		}

		log.Printf("Failed to ping DB (attempt %d): %v", i+1, err)
		time.Sleep(5 * time.Second) // 5秒待ってから再試行
	}

	log.Fatalln(err)

	return nil
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
