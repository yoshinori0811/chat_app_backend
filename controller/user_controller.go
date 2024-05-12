package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yoshinori0811/chat_app/config"
	"github.com/yoshinori0811/chat_app/model"
	"github.com/yoshinori0811/chat_app/usecase"
)

type IUserController interface {
	// UserControllerのメソッドを記載
	SignUP(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &UserController{uu}
}

// 〇TEST: POST以外の場合、"Method not allowed" が返ること
// 〇TEST: Content-typeが異なる場合、"Bad request"を返すこと
// 〇TEST: ユーザー名、メールアドレス、パスワードのいずれかがゼロ値の場合、"Bad request" を返すこと（レコードが作成されないこと）
// 〇TEST: すでに存在するユーザー名もしくはメールアドレスを受け取ったら"Internal server error" を返すこと（レコードが作成されないこと）
// 〇TEST: ステータスが201で返り、ユーザー名、メールアドレス、パスワード等が返ること
func (uc *UserController) SignUP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Star signup request")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq model.UserSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
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
	// CAUTION: エラーの内容によって条件分岐したい
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRes)
}

// 〇TEST: ステータス200で返ること[正常系]
// 〇TEST: cookieが生成されること[正常系]
// 〇TEST: POST以外でリクエストされたら "Method not allowed" を返す[異常系]
// 〇TEST: POSTデータのバインドに失敗したら "Bad request" を返す[異常系]
// 〇TEST: ユーザー名、パスワードのいずれかが存在しない場合 "Bad request" を返す[異常系]
// TEST: パスワードが異なる場合 "Internal server error" を返す[異常系]
// 〇TEST: uc.uu.Login()でエラーが返ったら "Internal server error" を返す（sessionIDの生成失敗など）[異常系]
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Star login request")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// POSTデータのバインド
	var userReq model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Println("userReq", userReq)
	if userReq.Email == "" || userReq.Password == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user := model.User{
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	fmt.Println("user", user)

	sessionID, err := uc.uu.Login(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = config.Config.ServerDomain
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

// TEST: ステータスが200で返ること
// TEST: cookieの有効期限が過ぎてること, 値が空文字であること
// TEST: GET以外 "Method not allowed" を返す
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
