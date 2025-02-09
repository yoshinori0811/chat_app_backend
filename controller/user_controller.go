package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/schema"

	"github.com/yoshinori0811/chat_app_backend/config"
	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/usecase"
)

type UserControllerInterface interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	SearchUsers(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	uu usecase.UserUsecaseInterface
	fu usecase.FriendUsecaseInterface
}

func NewUserController(uu usecase.UserUsecaseInterface, fu usecase.FriendUsecaseInterface) UserControllerInterface {
	return &UserController{uu, fu}
}

var decoder = schema.NewDecoder()

func (uc *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	// if !checkPOSTMethod(w, r) {
	// 	return
	// }

	userReq, err := bindJSON[model.UserSignUpRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userReq)
	if userReq.UserName == "" || userReq.Email == "" || userReq.Password == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:     userReq.UserName,
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRes)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// if !checkPOSTMethod(w, r) {
	// 	return
	// }

	userReq, err := bindJSON[model.UserLoginRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	if userReq.Email == "" || userReq.Password == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user := model.User{
		Email:    userReq.Email,
		Password: userReq.Password,
	}
	session, err := uc.uu.Login(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = session.SessionToken
	cookie.Expires = session.ExpiredAt
	cookie.Path = "/"
	// MEMO: Docker化したアプリをローカルで実行する場合Domainを""とする
	// cookie.Domain = ""
	cookie.Domain = config.Config.ServerDomain
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode

	fmt.Println("session.SessionToken:", session.SessionToken)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// if !checkPOSTMethod(w, r) {
	// 	return
	// }

	// リクエストからCookieを取得
	cookies := r.Cookies()
	// 取得したCookieを表示
	for _, cookie := range cookies {
		fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
	}

	sessionID, err := r.Cookie("session")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	http.SetCookie(w, cookie)

	if err := uc.uu.Logout(sessionID.Value); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDContextKey).(uint)
	res, err := uc.uu.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("ユーザー取得ロジック")
	json.NewEncoder(w).Encode(res)
}

func (uc *UserController) SearchUsers(w http.ResponseWriter, r *http.Request) {
	// GETデータのバインド
	userReq, err := bindQueryParams[model.UserSearchRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("throw bindQueryParams")

	userID := r.Context().Value(model.UserIDContextKey).(uint)
	res, err := uc.uu.SearchUsers(userReq.Query, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

// func checkPOSTMethod(w http.ResponseWriter, r *http.Request) bool {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return false
// 	}
// 	return true
// }

// func checkGETMethod(w http.ResponseWriter, r *http.Request) bool {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return false
// 	}
// 	return true
// }

func bindJSON[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil, err
	}
	return &data, nil
}

func bindQueryParams[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var data T
	if err := decoder.Decode(&data, r.URL.Query()); err != nil {
		fmt.Println(err)
		fmt.Println("bad request:", err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil, err
	}
	return &data, nil
}
