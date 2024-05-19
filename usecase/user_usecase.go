package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yoshinori0811/chat_app/model"
	"github.com/yoshinori0811/chat_app/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (model.Session, error)
	Logout(sessionToken string) error
}

type userUsecase struct {
	ur repository.IUserRepository
	sr repository.ISessionRepository
}

func NewUserUsecase(ur repository.IUserRepository, sr repository.ISessionRepository) IUserUsecase {
	return &userUsecase{
		ur,
		sr,
	}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	if storedUser.Email == user.Email || storedUser.Name == user.Name {
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
		ID:    newUser.ID,
		UUID:  newUser.UUID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}
	fmt.Println("created record: ", newUser)
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (model.Session, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}
	// パスワード検証
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}
	// セッションの生成、保存
	sessionToken, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}
	newSession := model.Session{
		SessionToken: sessionToken.String(),
		ExpiredAt:    time.Now().Add(24 * time.Hour),
	}
	if err := uu.sr.CreateSession(&newSession, storedUser.ID); err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}

	return newSession, nil
}

func (uu *userUsecase) Logout(sessionID string) error {
	if err := uu.sr.DeleteSession(sessionID); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
