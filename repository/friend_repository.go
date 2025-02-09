package repository

import (
	"github.com/yoshinori0811/chat_app_backend/model"
	"gorm.io/gorm"
)

type FriendRepositoryInterface interface {
	GetFriendsByUserID(userID uint, roomType uint) ([]model.FriendResponse, error)
	InsertFriendPair(sender model.Friend, receiver model.Friend) error
	GetFriendsWithMessagesDesc(userID uint, roomType uint) ([]model.FriendResponse, error)
}

type FriendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) FriendRepositoryInterface {
	return &FriendRepository{db}
}

func (fr FriendRepository) GetFriendsByUserID(userID uint, roomType uint) ([]model.FriendResponse, error) {
	var friends []model.FriendResponse
	sql := `SELECT u.name AS name, r.uuid AS room_uuid
		FROM rooms AS r
		JOIN (
			SELECT rm1.room_id, rm2.user_id
			FROM room_members rm1
			JOIN room_members rm2 ON rm1.room_id = rm2.room_id
			WHERE rm1.user_id = ?
			AND rm2.user_id <> ?
			AND rm1.room_id IN (
				SELECT room_id
				FROM room_members
				GROUP BY room_id
				HAVING COUNT(*) = 2
			)
		) AS rm
		ON r.id = rm.room_id
		LEFT JOIN users AS u
		ON rm.user_id = u.id
		WHERE r.type = ?
		ORDER BY r.last_message_at DESC`

	if err := fr.db.Raw(sql, userID, userID, roomType).Scan(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}

func (fr FriendRepository) GetFriendsWithMessagesDesc(userID uint, roomType uint) ([]model.FriendResponse, error) {
	var friends []model.FriendResponse
	sql := `SELECT u.name AS name, r.uuid AS room_uuid
		FROM rooms AS r
		JOIN (
			SELECT rm1.room_id, rm2.user_id
			FROM room_members rm1
			JOIN room_members rm2 ON rm1.room_id = rm2.room_id
			WHERE rm1.user_id = ?
			AND rm2.user_id <> ?
			AND rm1.room_id IN (
				SELECT room_id
				FROM room_members
				GROUP BY room_id
				HAVING COUNT(*) = 2
			)
		) AS rm
		ON r.id = rm.room_id
		LEFT JOIN users AS u
		ON rm.user_id = u.id
		WHERE r.last_message_at IS NOT NULL
		AND r.type = ?
		ORDER BY r.last_message_at DESC`
	if err := fr.db.Raw(sql, userID, userID, roomType).Scan(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}

func (fr FriendRepository) InsertFriendPair(sender model.Friend, receiver model.Friend) error {
	sql := `INSERT INTO friends (user_id, friend_id) VALUES (?, ?)`
	tx := fr.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Exec(sql, sender.UserID, sender.FriendID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec(sql, receiver.UserID, receiver.FriendID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
