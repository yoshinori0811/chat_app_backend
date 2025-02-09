package repository

import (
	"fmt"

	"github.com/yoshinori0811/chat_app_backend/model"
	"gorm.io/gorm"
)

type RoomRepositoryInterface interface {
	Insert(room *model.Room, tx *gorm.DB) error
	GetByUUID(room *model.Room) error
	GetUUIDAndNameByRoomMemberUserID(userID uint, roomType uint) ([]model.GetRoomsResponse, error)
	DeleteByRoomUUID(room *model.Room) error
	UpdateLastMessageAtByRoomUUID(room *model.Room) error
}

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepositoryInterface {
	return &RoomRepository{db}
}

func (rr *RoomRepository) Insert(room *model.Room, tx *gorm.DB) error {
	fmt.Println("before Insert Room: ", room)
	if err := tx.Select("uuid", "name", "admin_user_id", "type").Create(&room).Error; err != nil {
		return err
	}
	fmt.Println("after Insert Room: ", room)
	return nil
}

func (rr *RoomRepository) GetByUUID(room *model.Room) error {
	sql := `SELECT * FROM rooms WHERE uuid = ?`
	if err := rr.db.Raw(sql, room.UUID).Scan(&room).Error; err != nil {
		return err
	}
	return nil
}

func (rr RoomRepository) GetUUIDAndNameByRoomMemberUserID(userID uint, roomType uint) ([]model.GetRoomsResponse, error) {
	var rooms []model.GetRoomsResponse
	sql := `SELECT r.uuid AS uuid, IFNULL(r.name, "") AS name
		FROM rooms AS r
		LEFT JOIN room_members AS rm
		ON r.id = rm.room_id
		WHERE rm.user_id = ?
		AND r.type = ?
		ORDER BY r.created_at DESC
		LIMIT 20`

	rows, err := rr.db.Raw(sql, userID, roomType).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room model.GetRoomsResponse
		if err := rows.Scan(&room.UUID, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil

}

func (rr RoomRepository) DeleteByRoomUUID(room *model.Room) error {
	if err := rr.db.Where("uuid = ?", room.UUID).Delete(room).Error; err != nil {
		return err
	}
	return nil
}

func (rr RoomRepository) UpdateLastMessageAtByRoomUUID(room *model.Room) error {
	if err := rr.db.Select("last_message_at").Where("id = ?", room.ID).Updates(room).Error; err != nil {
		return err
	}
	return nil
}
