package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"strconv"

	pb "github.com/yoshinori0811/chat_app_backend/pb"
	server "github.com/yoshinori0811/chat_app_backend/server/interceptor"
	"github.com/yoshinori0811/chat_app_backend/server/service"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/yoshinori0811/chat_app_backend/config"
	"github.com/yoshinori0811/chat_app_backend/controller"
	"github.com/yoshinori0811/chat_app_backend/db"
	"github.com/yoshinori0811/chat_app_backend/middleware"
	"github.com/yoshinori0811/chat_app_backend/repository"
	"github.com/yoshinori0811/chat_app_backend/router"
	"github.com/yoshinori0811/chat_app_backend/usecase"
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

	var messageService *service.MessageServiceServer
	var grpcServer *grpc.Server

	if config.Config.AppEnv == "production" {
		creds, err := credentials.NewServerTLSFromFile(config.Config.CertFile, config.Config.KeyFile)
		if err != nil {
			log.Fatalf("Failed to load TLS credentials: %v\n", err)
		}
		messageService = service.NewMessageServiceServer(roomUsecase)
		interceptor := server.NewInterceptor(sessionUsecase)
		grpcServer = grpc.NewServer(
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time:    30 * time.Second,
				Timeout: 10 * time.Second,
			}),
			grpc.Creds(creds),
			grpc.UnaryInterceptor(interceptor.UnarySessionInterceptor),
			grpc.StreamInterceptor(interceptor.ServerStreamSessionInterceptor),
		)
	} else {
		messageService = service.NewMessageServiceServer(roomUsecase)
		interceptor := server.NewInterceptor(sessionUsecase)
		grpcServer = grpc.NewServer(
			grpc.UnaryInterceptor(interceptor.UnarySessionInterceptor),
			grpc.StreamInterceptor(interceptor.ServerStreamSessionInterceptor),
		)
	}

	pb.RegisterMessageServiceServer(grpcServer, messageService)
	lis, err := net.Listen("tcp", ":"+config.Config.ServerGrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	go func() {
		if config.Config.AppEnv == "production" {
			if err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Config.ServerPort), config.Config.CertFile, config.Config.KeyFile, nil); err != nil {
				log.Fatalf("Failed to serve HTTPS: %v\n", err)
			}
			log.Println(strconv.Itoa(config.Config.ServerPort))
		} else {
			if err := http.ListenAndServe(":"+strconv.Itoa(config.Config.ServerPort), nil); err != nil {
				log.Fatalf("Failed to serve HTTP: %v\n", err)
			}
			log.Println(strconv.Itoa(config.Config.ServerPort))
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v\n", err)
	}
}
