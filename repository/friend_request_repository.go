package repository

import (
	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/model/enum"
	"gorm.io/gorm"
)

type FriendRequestRepositoryInterface interface {
	Insert(senderID uint, receiverID uint) error
	GetFriendRequestsByReceiverID(receiverID uint) ([]model.FriendRequestListResponse, error)
	FindByReceiverID(senderID uint, receiverID uint) (*model.FriendRequest, error)
	UpdateStatusByreceiverIDAndSenderID(friendRequest *model.FriendRequest) error
}

type FriendRequestRepository struct {
	db *gorm.DB
}

func NewFriendRequestRepository(db *gorm.DB) FriendRequestRepositoryInterface {
	return &FriendRequestRepository{db}
}

func (frr FriendRequestRepository) Insert(senderID uint, receiverID uint) error {
	sql := `INSERT INTO friend_requests (sender_id, receiver_id) VALUES (?, ?)`
	if err := frr.db.Exec(sql, senderID, receiverID).Error; err != nil {
		return err
	}
	return nil
}

func (frr FriendRequestRepository) GetFriendRequestsByReceiverID(receiverID uint) ([]model.FriendRequestListResponse, error) {
	friendRequests := []model.FriendRequestListResponse{}
	sql := `SELECT fr.status AS status, u.name AS user_name FROM friend_requests AS fr LEFT JOIN users AS u ON fr.sender_id = u.id WHERE receiver_id = ? AND status = ?`
	if err := frr.db.Raw(sql, receiverID, enum.Pending).Scan(&friendRequests).Error; err != nil {
		return nil, err
	}
	return friendRequests, nil
}

func (frr FriendRequestRepository) FindByReceiverID(senderID uint, receiverID uint) (*model.FriendRequest, error) {
	var friendRequest model.FriendRequest
	sql := `SELECT * FROM friend_requests WHERE sender_id = ? AND receiver_id = ?`
	if err := frr.db.Raw(sql, senderID, receiverID).First(&friendRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &friendRequest, nil
}

func (frr FriendRequestRepository) UpdateStatusByreceiverIDAndSenderID(friendRequest *model.FriendRequest) error {
	sql := `UPDATE friend_requests SET status = ? WHERE receiver_id = ? AND sender_id = ?`
	if err := frr.db.Exec(sql, friendRequest.Status, friendRequest.ReceiverID, friendRequest.SenderID).Error; err != nil {
		return err
	}
	return nil
}
