package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/usecase"
)

type FriendControllerInterface interface {
	SendFriendRequest(w http.ResponseWriter, r *http.Request)
	GetFriendRequestList(w http.ResponseWriter, r *http.Request)
	GetFriends(w http.ResponseWriter, r *http.Request)
	GetDmList(w http.ResponseWriter, r *http.Request)
	AcceptFriendRequest(w http.ResponseWriter, r *http.Request)
	RejectFriendRequest(w http.ResponseWriter, r *http.Request)
}

type FriendController struct {
	fu usecase.FriendUsecaseInterface
	ru usecase.RoomUsecaseInterface
}

func NewFriendController(fu usecase.FriendUsecaseInterface, ru usecase.RoomUsecaseInterface) FriendControllerInterface {
	return &FriendController{fu, ru}
}

func (fc *FriendController) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	reqBody, err := bindJSON[model.FriendRequestRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	senderID := r.Context().Value(model.UserIDContextKey).(uint)
	if err := fc.fu.SendFriendRequest(senderID, reqBody.UserName); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (fc *FriendController) GetFriendRequestList(w http.ResponseWriter, r *http.Request) {
	receiverID := r.Context().Value(model.UserIDContextKey).(uint)

	res, err := fc.fu.GetFriendRequestList(receiverID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (fc *FriendController) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDContextKey).(uint)

	res, err := fc.fu.GetFriends(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

func (fc *FriendController) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	reqBody, err := bindJSON[model.FriendRequestRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	receiverID := r.Context().Value(model.UserIDContextKey).(uint)
	senderID, tx, err := fc.fu.AcceptFriendRequest(receiverID, reqBody.UserName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := fc.ru.CreateDMRoom(receiverID, senderID, tx); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (fc *FriendController) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	reqBody, err := bindJSON[model.FriendRequestRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	receiverID := r.Context().Value(model.UserIDContextKey).(uint)
	if err := fc.fu.RejectFriendRequest(receiverID, reqBody.UserName); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (fc *FriendController) GetDmList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDContextKey).(uint)

	res, err := fc.fu.GetFriendsWithMessagesDesc(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}
