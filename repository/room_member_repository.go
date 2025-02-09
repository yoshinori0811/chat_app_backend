package repository

import (
	"github.com/yoshinori0811/chat_app_backend/model"
	"gorm.io/gorm"
)

type RoomMemberRepositoryInterface interface {
	Insert(members []model.RoomMember, tx *gorm.DB) error
	GetRoomMemberNamesByRoomID(roomID uint) ([]string, error)
	DeleteByRoomIDAndUserID(member *model.RoomMember) error
}

type RoomMemberRepository struct {
	db *gorm.DB
}

func NewRoomMemberRepository(db *gorm.DB) RoomMemberRepositoryInterface {
	return &RoomMemberRepository{db}
}

// parameters:
// -tx: nilの場合、トランザクションを行わない
func (rr *RoomMemberRepository) Insert(members []model.RoomMember, tx *gorm.DB) error {
	if tx == nil {
		if err := rr.db.Select("room_id", "user_id").Create(&members).Error; err != nil {
			return err
		}
		return nil
	}

	if err := tx.Select("room_id", "user_id").Create(&members).Error; err != nil {
		return err
	}
	return nil
}

func (rr RoomMemberRepository) GetRoomMemberNamesByRoomID(roomID uint) ([]string, error) {
	var roomMembers []string
	sql := `SELECT u.name AS name
		FROM room_members AS rm
		LEFT JOIN users AS u
		ON rm.user_id = u.id
		WHERE room_id = ?`
	if err := rr.db.Raw(sql, roomID).Scan(&roomMembers).Error; err != nil {
		return nil, err
	}
	return roomMembers, nil
}

func (rr RoomMemberRepository) DeleteByRoomIDAndUserID(member *model.RoomMember) error {
	if err := rr.db.Where("room_id = ? AND user_id = ?", member.RoomID, member.UserID).Delete(member).Error; err != nil {
		return err
	}
	return nil
}
