package main

import (
	"net/http"
	"strconv"

	"github.com/yoshinori0811/chat_app/config"
	"github.com/yoshinori0811/chat_app/controller"
	"github.com/yoshinori0811/chat_app/db"
	"github.com/yoshinori0811/chat_app/repository"
	"github.com/yoshinori0811/chat_app/router"
	"github.com/yoshinori0811/chat_app/usecase"
)

func main() {
	db := db.NewDB()
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, sessionRepository)
	userController := controller.NewUserController(userUsecase)

	router.NewRouter(userController)
	http.ListenAndServe(":"+strconv.Itoa(config.Config.ServerPort), nil)
}
