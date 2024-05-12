package usecase

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yoshinori0811/chat_app/model"
	"github.com/yoshinori0811/chat_app/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	if storedUser.Email == user.Email || storedUser.Name == user.Name {
		fmt.Println("storedUser: ", &storedUser)
		fmt.Println("user: ", &user)
		err := errors.New("email or name already exists")
		fmt.Println(err)
		return model.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	// UUIDを生成する
	uuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	newUser := model.User{UUID: uuid.String(), Name: user.Name, Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:        newUser.ID,
		UUID:      newUser.UUID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		SessionID: newUser.SessionID,
	}
	fmt.Println("created record: ", newUser)
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	fmt.Println(user.Email)
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		fmt.Println(err)
		return "", err
	}
	// パスワード検証
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	fmt.Println("hash: ", encodedHash)
	fmt.Println("stored hash: ", storedUser.Password)
	fmt.Println("stored user: ", storedUser)

	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// セッションの生成、保存
	sessionID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if err := uu.ur.UpdateSession(&user, sessionID.String()); err != nil {
		fmt.Println(err)
		return "", err
	}

	return sessionID.String(), nil
}
