package model

import "time"

type Session struct {
	ID           uint      `json:"id" gorm:"primaryKey;"`
	UserID       uint      `json:"user_id" gorm:"not null;"`
	User         User      `json:"user" gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	SessionToken string    `json:"session_token" gorm:"unique;not null;"`
	ExpiredAt    time.Time `json:"expired_at" gorm:"not null;"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type ContextKey string

const (
	UserIDContextKey = ContextKey("UserID")
)
