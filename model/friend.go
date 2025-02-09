package model

import (
	"time"

	"github.com/yoshinori0811/chat_app_backend/model/enum"
)

type FriendRequest struct {
	ID         uint                     `json:"id" gorm:"primaryKey;"`
	SenderID   uint                     `json:"sender_id" gorm:"not null;uniqueIndex:idx_sender_receiver;"`
	ReceiverID uint                     `json:"receiver_id" gorm:"not null;uniqueIndex:idx_sender_receiver;"`
	Status     enum.FriendRequestStatus `json:"status" gorm:"not null;type:enum('accept','pending','reject');default:'pending';"`
	CreatedAt  time.Time                `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt  time.Time                `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt  time.Time                `json:"deleted_at"`
	Sender     User                     `json:"sender" gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Receiver   User                     `json:"receiver" gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type Friend struct {
	ID        uint      `json:"id" gorm:"primaryKey;"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_user_id_friend_id"`
	FriendID  uint      `json:"friend_id" gorm:"uniqueIndex:idx_user_id_friend_id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt time.Time `json:"deleted_at"`
	User      User      `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Friend    User      `json:"Friend" gorm:"foreignKey:FriendID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type FriendRequestRequest struct {
	UserName string `json:"user_name"`
}

type FriendRequestListResponse struct {
	UserName string                   `json:"user_name"`
	Status   enum.FriendRequestStatus `json:"status"`
}

type FriendResponse struct {
	Name     string `json:"name"`
	RoomUUID string `json:"room_uuid"`
}
