package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseInterface interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (model.Session, error)
	Logout(sessionToken string) error
	isEmailExists(email string) error
	SearchUsers(name string, userID uint) ([]model.UserSearchResponse, error)
	GetUser(id uint) (model.UserInfo, error)
}

type userUsecase struct {
	ur repository.UserRepositoryInterface
	sr repository.SessionRepositoryInterface
}

func NewUserUsecase(ur repository.UserRepositoryInterface, sr repository.SessionRepositoryInterface) UserUsecaseInterface {
	return &userUsecase{
		ur,
		sr,
	}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.isEmailExists(user.Email); err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	newUser := model.User{UUID: uuid.String(), Name: user.Name, Email: user.Email, Password: string(hash)}
	if err := uu.ur.Insert(&newUser); err != nil {
		fmt.Println(err)
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:    newUser.ID,
		UUID:  newUser.UUID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (model.Session, error) {
	storedUser := model.User{}
	if err := uu.ur.GetByEmail(&storedUser, user.Email); err != nil {
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
	if err := uu.sr.Insert(&newSession, storedUser.ID); err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}

	return newSession, nil
}

func (uu *userUsecase) Logout(sessionID string) error {
	if err := uu.sr.DeleteBySessionToken(sessionID); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (uu userUsecase) isEmailExists(email string) error {
	isUser, err := uu.ur.ExistsUserByEmail(email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if isUser {
		err := errors.New("email already exists")
		return err
	}
	return nil
}

func (uu *userUsecase) SearchUsers(name string, userID uint) ([]model.UserSearchResponse, error) {
	users, err := uu.ur.GetByName(name, userID) // []model.Userが返る
	if err != nil {
		return []model.UserSearchResponse{}, err
	}
	var userSeachRes []model.UserSearchResponse
	for _, v := range users {
		// userSeachResにマッピングする
		userSeachRes = append(userSeachRes, model.UserSearchResponse{
			Name: v.Name,
			UUID: v.UUID,
		})
	}
	return userSeachRes, nil
}

func (uu userUsecase) GetUser(id uint) (model.UserInfo, error) {
	user := model.User{
		ID: id,
	}
	if err := uu.ur.GetUserByID(&user); err != nil {
		fmt.Println(err)
		return model.UserInfo{}, err
	}
	res := model.UserInfo{
		Name: user.Name,
	}
	return res, nil
}
