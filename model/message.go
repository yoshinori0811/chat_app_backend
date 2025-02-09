package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey;"`
	UUID      string         `json:"uuid" gorm:"not null;unique"`
	UserID    uint           `json:"user_id" gorm:"not null;"`
	RoomID    uint           `json:"room_id" gorm:"not null;"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Room      Room           `json:"room" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	User      User           `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type MessageInfo struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"` // MEMO: Message.CreatedAtが格納される
	User      UserInfo  `json:"user"`
}

type BroadcastMessage struct {
	Type        string      `json:"type"`
	MessageInfo MessageInfo `json:"message_info"`
}

type MessageCreateRequest struct {
	Content string `json:"content"`
}

type MessageUpdateRequest struct {
	Content string `json:"content"`
}
