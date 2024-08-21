package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;"`
	UUID      string    `json:"uuid" gorm:"unique;"`
	Name      string    `json:"name" gorm:"unique;not null;"`
	Email     string    `json:"email" gorm:"unique; not null;"`
	Password  string    `json:"password" gorm:"not null;"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt time.Time `json:"deleted_at"`
	Session   []Session `json:"session" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type UserChat struct {
	ID            uint `json:"id" gorm:"primaryKey;"`
	UserID        uint
	FriendUserID  uint
	LastMessageAt time.Time
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt     time.Time `json:"deleted_at"`
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
	ID    uint   `json:"id"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserSearchRequest struct {
	Query string `json:"id" schema:"query,required"`
}

type UserSearchResponse struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type UserInfo struct {
	Name string `json:"name"`
}
