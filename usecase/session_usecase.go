package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/repository"
)

type SessionUsecaseInterface interface {
	ValidateSession(sessionToken string) (model.Session, error)
}

type SessionUsecase struct {
	sr repository.SessionRepositoryInterface
}

func NewSessionUsecase(sr repository.SessionRepositoryInterface) SessionUsecaseInterface {
	return &SessionUsecase{sr}
}

func (su *SessionUsecase) ValidateSession(sessionToken string) (model.Session, error) {
	session := model.Session{
		SessionToken: sessionToken,
	}
	if err := su.sr.GetBySessionToken(&session); err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}
	if session.ExpiredAt.Before(time.Now()) {
		return model.Session{}, errors.New("session expired")
	}
	return session, nil
}
