package router

import (
	"net/http"

	"github.com/yoshinori0811/chat_app/config"
	"github.com/yoshinori0811/chat_app/controller"
)

func NewRouter(uc controller.IUserController) {
	http.HandleFunc("/signup", corsMiddleware(uc.SignUp))
	http.HandleFunc("/login", corsMiddleware(uc.Login))
	http.HandleFunc("/logout", corsMiddleware(uc.Logout))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.Config.FEUrl)
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,X-CSRF-Header,Accept,Access-Control-AllowHeaders,")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
