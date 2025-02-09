package server

import (
	"context"
	"fmt"

	"net/http"

	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type serverStreamWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

type Interceptor struct {
	su usecase.SessionUsecaseInterface
}

func NewInterceptor(su usecase.SessionUsecaseInterface) *Interceptor {
	return &Interceptor{
		su,
	}
}

func (i *Interceptor) ServerStreamSessionInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	fmt.Printf("Received method: %s\n", info.FullMethod)
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		fmt.Println("SessionStreamInterceptor !ok:", md)
		return status.Errorf(codes.Unauthenticated, "No metadata in context")
	}

	cookies := md["cookie"]
	var sessionToken string

	for _, cookie := range cookies {
		httpCookies := &http.Request{Header: http.Header{"Cookie": []string{cookie}}}
		if c, err := httpCookies.Cookie("session"); err == nil {
			sessionToken = c.Value
			break
		}
	}

	if sessionToken == "" {
		fmt.Println("SessionStreamInterceptor sessioinToken is null")
		return status.Errorf(codes.Unauthenticated, "Sesion ID is missing or invalid")
	}

	session, err := i.su.ValidateSession(sessionToken)
	if err != nil {
		fmt.Println("SessionStreamInterceptor  ValidateSession:", err)
		status.Errorf(codes.Unauthenticated, "Sesion ID is missing or invalid")
	}

	newCtx := context.WithValue(ss.Context(), model.UserIDContextKey, session.UserID)
	wrappedStream := &serverStreamWrapper{ServerStream: ss, ctx: newCtx}

	return handler(srv, wrappedStream)
}

func (i *Interceptor) UnarySessionInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("Received method: %s\n", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("SessionStreamInterceptor !ok:", md)
		return nil, status.Errorf(codes.Unauthenticated, "No metadata in context")
	}

	cookies := md["cookie"]
	var sessionToken string
	for _, cookie := range cookies {
		httpCookies := &http.Request{Header: http.Header{"Cookie": []string{cookie}}}
		if c, err := httpCookies.Cookie("session"); err == nil {
			sessionToken = c.Value
			break
		}
	}

	if sessionToken == "" {
		fmt.Println("SessionStreamInterceptor sessioinToken is null")
		return nil, status.Errorf(codes.Unauthenticated, "Sesion ID is missing or invalid")
	}

	session, err := i.su.ValidateSession(sessionToken)
	if err != nil {
		fmt.Println("SessionStreamInterceptor  ValidateSession:", err)
		status.Errorf(codes.Unauthenticated, "Sesion ID is missing or invalid")
	}

	newCtx := context.WithValue(ctx, model.UserIDContextKey, session.UserID)
	return handler(newCtx, req)
}

func (s *serverStreamWrapper) Context() context.Context {
	ctx := s.ctx
	return ctx
}
