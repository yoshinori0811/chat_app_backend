package main

import (
	"fmt"

	"github.com/yoshinori0811/chat_app/db"
	"github.com/yoshinori0811/chat_app/model"
)

func main() {
	dbCon := db.NewDB()
	defer db.CloseDB(dbCon)
	dbCon.AutoMigrate(&model.User{}, &model.Session{})
	fmt.Println("Successfully Migrated")
}
