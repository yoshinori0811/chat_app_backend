package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yoshinori0811/chat_app_backend/config"
	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/usecase"
)

type MiddlewareInterface interface {
	CorsMiddleware(next http.Handler) http.HandlerFunc
	AuthMiddleware(next http.Handler) http.HandlerFunc
}

type Middleware struct {
	su usecase.SessionUsecaseInterface
}

func NewMiddleware(su usecase.SessionUsecaseInterface) MiddlewareInterface {
	return &Middleware{su}
}

func (m *Middleware) CorsMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request: %s %s %s\n", r.URL.String(), r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", config.Config.FEUrl)
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,X-CSRF-Header,Accept,Access-Control-AllowHeaders")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		session, err := m.su.ValidateSession(sessionCookie.Value)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), model.UserIDContextKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
