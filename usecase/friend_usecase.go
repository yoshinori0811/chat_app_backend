package usecase

import (
	"errors"
	"fmt"

	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/model/enum"
	"github.com/yoshinori0811/chat_app_backend/repository"
	"gorm.io/gorm"
)

type FriendUsecaseInterface interface {
	SendFriendRequest(senderID uint, receiverName string) error
	GetFriends(userID uint) ([]model.FriendResponse, error)
	GetFriendsWithMessagesDesc(userID uint) ([]model.FriendResponse, error)
	GetFriendRequestList(receiverID uint) ([]model.FriendRequestListResponse, error)
	AcceptFriendRequest(receiverID uint, senderName string) (uint, *gorm.DB, error)
	RejectFriendRequest(receiverID uint, senderName string) error
}

type FriendUsecase struct {
	ur  repository.UserRepositoryInterface
	frr repository.FriendRequestRepositoryInterface
	fr  repository.FriendRepositoryInterface
	db  *gorm.DB
}

func NewFriendUsecase(ur repository.UserRepositoryInterface, frr repository.FriendRequestRepositoryInterface, fr repository.FriendRepositoryInterface, db *gorm.DB) FriendUsecaseInterface {
	return &FriendUsecase{ur, frr, fr, db}
}

func (fu *FriendUsecase) SendFriendRequest(senderID uint, receiverName string) error {
	receiverID, err := fu.ur.GetUserIDByName(receiverName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	friendRequest, err := fu.frr.FindByReceiverID(senderID, receiverID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if friendRequest == nil {
		if err := fu.frr.Insert(senderID, receiverID); err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		if friendRequest.Status != enum.Reject {
			return errors.New("friend request status is already accept or pending")
		}
		friendRequest.Status = enum.Pending
		if err := fu.frr.UpdateStatusByreceiverIDAndSenderID(friendRequest); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (fu *FriendUsecase) GetFriendRequestList(receiverID uint) ([]model.FriendRequestListResponse, error) {
	friendRequests, err := fu.frr.GetFriendRequestsByReceiverID(receiverID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return friendRequests, nil
}

func (fu *FriendUsecase) GetFriends(userID uint) ([]model.FriendResponse, error) {
	friends, err := fu.fr.GetFriendsByUserID(userID, 1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return friends, nil
}

func (fu *FriendUsecase) GetFriendsWithMessagesDesc(userID uint) ([]model.FriendResponse, error) {
	friends, err := fu.fr.GetFriendsWithMessagesDesc(userID, 1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return friends, nil
}

func (fu *FriendUsecase) AcceptFriendRequest(receiverID uint, senderName string) (uint, *gorm.DB, error) {
	senderID, err := fu.ur.GetUserIDByName(senderName)
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	friendRequest := model.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     enum.Accept,
	}

	tx := fu.db.Begin()
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return 0, nil, tx.Error
	}

	if err := fu.frr.UpdateStatusByreceiverIDAndSenderID(&friendRequest); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return 0, nil, err
	}

	sender := model.Friend{
		UserID:   friendRequest.SenderID,
		FriendID: friendRequest.ReceiverID,
	}
	receiver := model.Friend{
		UserID:   friendRequest.ReceiverID,
		FriendID: friendRequest.SenderID,
	}

	if err := fu.fr.InsertFriendPair(sender, receiver); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return 0, nil, err
	}
	return senderID, tx, nil
}

func (fu *FriendUsecase) RejectFriendRequest(receiverID uint, senderName string) error {
	senderID, err := fu.ur.GetUserIDByName(senderName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	friendRequest := model.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     enum.Reject,
	}

	if err := fu.frr.UpdateStatusByreceiverIDAndSenderID(&friendRequest); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
