package repository

import (
	"fmt"

	"github.com/yoshinori0811/chat_app_backend/model"
	"gorm.io/gorm"
)

type MessageRepositoryInterface interface {
	GetMessagesByRoomID(roomID uint, offset uint) ([]model.MessageInfo, error)
	GetByID(message *model.Message) error
	GetMessageByID(ID uint) (model.MessageInfo, error)
	GetByUUID(message *model.Message) error
	Insert(message *model.Message) error
	UpdateContentByUUID(message *model.Message) error
	DeleteByUUID(messageUUID string) error
}

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepositoryInterface {
	return MessageRepository{db}
}

func (mr MessageRepository) GetMessagesByRoomID(roomID uint, offset uint) ([]model.MessageInfo, error) {
	var messages []model.MessageInfo
	sql := `SELECT m.id AS id, m.uuid AS uuid, m.content AS content, m.created_at AS timestamp, u.name AS user_name
		FROM messages AS m
		LEFT JOIN users AS u
		ON m.user_id = u.id
		WHERE m.room_id = ?
		ORDER BY m.created_at DESC
		LIMIT 50
		OFFSET ?`
	rows, err := mr.db.Raw(sql, roomID, offset).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mInfo model.MessageInfo
		var uInfo model.UserInfo
		if err := rows.Scan(&mInfo.ID, &mInfo.UUID, &mInfo.Content, &mInfo.Timestamp, &uInfo.Name); err != nil {
			return nil, err
		}
		mInfo.User = uInfo
		messages = append(messages, mInfo)
	}
	return messages, nil
}

func (mr MessageRepository) GetByID(message *model.Message) error {
	sql := `SELECT * FROM messages WHERE ID = ?`
	if err := mr.db.Raw(sql, message.ID).First(message).Error; err != nil {
		return err
	}
	return nil
}

func (mr MessageRepository) GetMessageByID(ID uint) (model.MessageInfo, error) {
	var mInfo model.MessageInfo
	var uInfo model.UserInfo
	sql := `SELECT m.id AS id, m.uuid AS uuid, m.content AS content, m.created_at AS timestamp, u.name AS user_name
		FROM messages AS m
		LEFT JOIN users AS u
		ON m.user_id = u.id
		WHERE m.id = ?`
	row := mr.db.Raw(sql, ID).Row()

	if err := row.Scan(&mInfo.ID, &mInfo.UUID, &mInfo.Content, &mInfo.Timestamp, &uInfo.Name); err != nil {
		return model.MessageInfo{}, err
	}
	mInfo.User = uInfo
	return mInfo, nil
}

func (mr MessageRepository) GetByUUID(message *model.Message) error {
	sql := `SELECT m.id AS id, m.uuid AS uuid, m.content AS content, m.created_at AS created_at, u.name AS user_name
		FROM messages AS m
		LEFT JOIN users AS u
		ON m.user_id = u.id
		WHERE m.uuid = ?`
	if err := mr.db.Raw(sql, message.UUID).First(message).Error; err != nil {
		return err
	}

	return nil
}

func (mr MessageRepository) Insert(message *model.Message) error {
	if err := mr.db.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

func (mr MessageRepository) UpdateContentByUUID(message *model.Message) error {
	if err := mr.db.Select("content").Where("uuid = ?", message.UUID).Updates(message).Error; err != nil {
		fmt.Println("UpdateMessageByMessageUUID: ", err)
		return err
	}
	return nil
}

func (mr MessageRepository) DeleteByUUID(messageUUID string) error {
	sql := `DELETE FROM messages WHERE uuid = ?`
	if err := mr.db.Exec(sql, messageUUID).Error; err != nil {
		return err
	}
	return nil
}
