package main

import (
	"fmt"

	"github.com/yoshinori0811/chat_app/db"
	"github.com/yoshinori0811/chat_app/model"
)

func main() {
	dbCon := db.NewDB()
	defer db.CloseDB(dbCon)
	dbCon.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Friend{},
		&model.FriendRequest{},
		&model.Room{},
		&model.RoomMember{},
		&model.Message{},
	)
	fmt.Println("Successfully Migrated")
}
