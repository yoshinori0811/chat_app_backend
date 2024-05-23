package repository

import (
	"github.com/yoshinori0811/chat_app/model"
	"gorm.io/gorm"
)

type SessionRepositoryInterface interface {
	InsertSession(session *model.Session, userID uint) error
	// UpdateSession() error
	DeleteSession(sessionToken string) error
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepositoryInterface {
	return &SessionRepository{db}
}

func (sr SessionRepository) InsertSession(session *model.Session, userID uint) error {
	sql := `INSERT INTO sessions (user_id, session_token, expired_at) VALUES (?, ?, ?)`
	if err := sr.db.Exec(sql, userID, session.SessionToken, session.ExpiredAt).Error; err != nil {
		return err
	}
	return nil
}

func (sr SessionRepository) DeleteSession(sessionToken string) error {
	sql := `DELETE FROM sessions WHERE session_token = ?`
	if err := sr.db.Exec(sql, sessionToken).Error; err != nil {
		return err
	}
	return nil
}
