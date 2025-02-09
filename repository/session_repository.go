package repository

import (
	"time"

	"github.com/yoshinori0811/chat_app_backend/model"
	"gorm.io/gorm"
)

type SessionRepositoryInterface interface {
	Insert(session *model.Session, userID uint) error
	DeleteBySessionToken(sessionToken string) error
	GetBySessionToken(session *model.Session) error
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepositoryInterface {
	return &SessionRepository{db}
}

func (sr SessionRepository) Insert(session *model.Session, userID uint) error {
	sql := `INSERT INTO sessions (user_id, session_token, expired_at) VALUES (?, ?, ?)`
	if err := sr.db.Exec(sql, userID, session.SessionToken, session.ExpiredAt).Error; err != nil {
		return err
	}
	return nil
}

func (sr SessionRepository) DeleteBySessionToken(sessionToken string) error {
	sql := `DELETE FROM sessions WHERE session_token = ?`
	if err := sr.db.Exec(sql, sessionToken).Error; err != nil {
		return err
	}
	return nil
}

func (sr SessionRepository) GetBySessionToken(session *model.Session) error {
	now := time.Now()
	sql := `SELECT * FROM sessions WHERE session_token = ? AND expired_at > ?`
	if err := sr.db.Raw(sql, session.SessionToken, now).First(session).Error; err != nil {
		return err
	}
	return nil
}
