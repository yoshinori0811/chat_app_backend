package main

import (
	"fmt"

	"github.com/yoshinori0811/chat_app/db"
	"github.com/yoshinori0811/chat_app/model"
)

func main() {
	dbCon := db.NewDB()
	defer db.CloseDB(dbCon)
	// dbCon.AutoMigrate(
	// 	&model.User{},
	// 	&model.Session{},
	// 	&model.Friend{},
	// 	&model.FriendRequest{},
	// 	&model.Room{},
	// 	&model.RoomMember{},
	// 	&model.Message{},
	// )
	// dbCon.AutoMigrate(&model.User{}, &model.Session{})
	// dbCon.AutoMigrate(&model.Friend{}, &model.FriendRequest{})
	// dbCon.AutoMigrate(&model.Friend{}, &model.Room{}, &model.RoomMember{})
	// dbCon.AutoMigrate(&model.Message{})
	dbCon.AutoMigrate(&model.Room{})
	fmt.Println("Successfully Migrated")
}
