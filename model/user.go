package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)"`
}

type UserSignUpRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UUID      string `json:"uuid" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	SessionID string `json:"session_id"`
}
