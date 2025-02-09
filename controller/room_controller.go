package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/gorilla/websocket"

	"github.com/yoshinori0811/chat_app_backend/model"
	"github.com/yoshinori0811/chat_app_backend/usecase"
)

type RoomControllerInterface interface {
	CreateRoom(w http.ResponseWriter, r *http.Request)
	GetRoomChat(w http.ResponseWriter, r *http.Request)
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetRooms(w http.ResponseWriter, r *http.Request)
	InviteRoom(w http.ResponseWriter, r *http.Request)
	DeleteRoom(w http.ResponseWriter, r *http.Request)
	LeaveRoom(w http.ResponseWriter, r *http.Request)
	UpdateMessage(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
}

type RoomController struct {
	ru usecase.RoomUsecaseInterface
}

func NewRoomController(ru usecase.RoomUsecaseInterface) RoomControllerInterface {
	return &RoomController{
		ru,
	}
}

func (rc *RoomController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	reqBody, err := bindJSON[model.RoomCreateRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	userID := r.Context().Value(model.UserIDContextKey).(uint)
	res, err := rc.ru.CreateRoom(*reqBody, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (rc *RoomController) GetRoomChat(w http.ResponseWriter, r *http.Request) {
	roomUUID := r.PathValue("roomUUID")
	if roomUUID == "" {
		http.Error(w, "Room UUID is required", http.StatusBadRequest)
		return
	}
	fmt.Println("roomUUID", roomUUID)

	userID := r.Context().Value(model.UserIDContextKey).(uint)
	res, err := rc.ru.GetRoomMessages(roomUUID, userID, 0)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (rc *RoomController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	roomUUID := r.PathValue("roomUUID")
	if roomUUID == "" {
		http.Error(w, "Room UUID is required", http.StatusBadRequest)
		return
	}

	reqBody, err := bindJSON[model.MessageCreateRequest](w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	userID := r.Context().Value(model.UserIDContextKey).(uint)

	msg, err := rc.ru.CreateMessage(roomUUID, *reqBody, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rc.ru.SendMessageToRoomChannel(roomUUID, msg)
	w.WriteHeader(http.StatusOK)
}

func (rc RoomController) GetRooms(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDContextKey).(uint)
	res, err := rc.ru.GetRooms(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (rc RoomController) InviteRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDContextKey).(uint)
	roomUUID := r.PathValue("roomUUID")
	res, err := rc.ru.InviteRoom(userID, roomUUID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internarl server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (rc RoomController) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomUUID := r.PathValue("roomUUID")
	if err := rc.ru.DeleteRoom(roomUUID); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rc RoomController) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDContextKey).(uint)
	roomUUID := r.PathValue("roomUUID")
	if err := rc.ru.LeaveRoom(userID, roomUUID); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rc RoomController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	reqBody, err := bindJSON[model.MessageUpdateRequest](w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	messageUUID := r.PathValue("messageUUID")

	msg, err := rc.ru.UpdateMessage(messageUUID, reqBody.Content)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	roomUUID := r.PathValue("roomUUID")
	rc.ru.SendMessageToRoomChannel(roomUUID, msg)
	w.WriteHeader(http.StatusOK)
	fmt.Println("UpdateMessage success")
}

func (rc RoomController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageUUID := r.PathValue("messageUUID")
	msg, err := rc.ru.DeleteMessage(messageUUID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	roomUUID := r.PathValue("roomUUID")
	rc.ru.SendMessageToRoomChannel(roomUUID, msg)
	w.WriteHeader(http.StatusOK)
}
