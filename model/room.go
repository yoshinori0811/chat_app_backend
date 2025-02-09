package model

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/yoshinori0811/chat_app_backend/pb"
)

type Room struct {
	ID            uint      `json:"id" gorm:"primaryKey;"`
	UUID          string    `json:"uuid" gorm:"not null;unique;"`
	Name          string    `json:"name" gorm:"default:null;"`
	AdminUserID   uint      `json:"admin_user_id" gorm:"default:null;"`
	Type          uint      `json:"type"`
	LastMessageAt time.Time `json:"last_message_at" gorm:"default:null;"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt     time.Time `json:"deleted_at"`
}

type RoomMember struct {
	ID        uint      `json:"id" gorm:"primaryKey;"`
	RoomID    uint      `json:"room_id" gorm:"not null;uniqueIndex:idx_room_id_user_id;"` // MEMO: RoomIDとUserIDの組み合わせの重複禁止
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_room_id_user_id;"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);"`
	DeletedAt time.Time `json:"deleted_at"`
	Room      Room      `json:"room" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	User      User      `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type RoomCreateRequest struct {
	Name          string   `json:"name"`
	AdminUserName string   `json:"admin_user_name"` // MEMO: DMのRoomを作成する場合、adminUserは無し、Roomを作成する場合、作成者をadminUserとしている
	Members       []string `json:"members"`
}

type RoomCreateResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type GetRoomsResponse struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type RoomInfoResponse struct {
	Name     string        `json:"name"`
	UUID     string        `json:"uuid"`
	IsAdmin  bool          `json:"is_admin"`
	Members  []string      `json:"members"`
	Messages []MessageInfo `json:"messages"`
}

type ChatRoom struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan BroadcastMessage
	Mutex     sync.Mutex
}

type RoomChannels struct {
	ClientChannels []chan *pb.MessageResponse
	Mu             sync.Mutex
}

type RoomInviteResponse struct {
	UUID string `json:"uuid"`
}

type GetGetRoomChatRequest struct {
	Offset uint `json:"offset"`
}
