package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	// pb "github.com/yoshinori0811/chat_app/pb/proto"
	pb "github.com/yoshinori0811/chat_app/pb"
	server "github.com/yoshinori0811/chat_app/server/interceptor"
	"github.com/yoshinori0811/chat_app/server/service"
	"google.golang.org/grpc"

	"github.com/yoshinori0811/chat_app/config"
	"github.com/yoshinori0811/chat_app/controller"
	"github.com/yoshinori0811/chat_app/db"
	"github.com/yoshinori0811/chat_app/middleware"
	"github.com/yoshinori0811/chat_app/repository"
	"github.com/yoshinori0811/chat_app/router"
	"github.com/yoshinori0811/chat_app/usecase"
)

func main() {
	db := db.NewDB()
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	friendRequestRepository := repository.NewFriendRequestRepository(db)
	friendRepository := repository.NewFriendRepository(db)
	roomRepository := repository.NewRoomRepository(db)
	roomMemberRepository := repository.NewRoomMemberRepository(db)
	messageRepository := repository.NewMessageRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, sessionRepository)
	friendUsecase := usecase.NewFriendUsecase(userRepository, friendRequestRepository, friendRepository, db)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepository, roomMemberRepository, userRepository, friendRepository, messageRepository, db)

	userController := controller.NewUserController(userUsecase, friendUsecase)
	friendController := controller.NewFriendController(friendUsecase, roomUsecase)
	roomController := controller.NewRoomController(roomUsecase)

	middleware := middleware.NewMiddleware(sessionUsecase)

	router.NewRouter(middleware, userController, friendController, roomController)
	// http.ListenAndServe(":"+strconv.Itoa(config.Config.ServerPort), nil)
	// fmt.Println(strconv.Itoa(config.Config.ServerPort))

	messageService := service.NewMessageServiceServer(roomUsecase)
	interceptor := server.NewInterceptor(sessionUsecase)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnarySessionInterceptor),
		grpc.StreamInterceptor(interceptor.ServerStreamSessionInterceptor),
	)

	pb.RegisterMessageServiceServer(grpcServer, messageService)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		return
	}

	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(config.Config.ServerPort), nil); err != nil {
			fmt.Printf("Failed to serve HTTP: %v\n", err)
		}
		fmt.Println(strconv.Itoa(config.Config.ServerPort))
	}()

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("Failed to serve gRPC: %v\n", err)
	}
}
