package usecase

import (
	"fmt"
	"sort"

	"github.com/rs/xid"
	"github.com/yoshinori0811/chat_app_backend/model"
	pb "github.com/yoshinori0811/chat_app_backend/pb"
	"github.com/yoshinori0811/chat_app_backend/repository"
	"gorm.io/gorm"
)

type RoomUsecaseInterface interface {
	CreateDMRoom(userID uint, receiverID uint, tx *gorm.DB) error
	CreateRoom(req model.RoomCreateRequest, userID uint) (model.RoomCreateResponse, error)
	GetRoomMessages(uuid string, userID uint, offset uint) (model.RoomInfoResponse, error)
	CreateMessage(roomUUID string, req model.MessageCreateRequest, userID uint) (model.BroadcastMessage, error)
	GetRooms(userID uint) ([]model.GetRoomsResponse, error)
	InviteRoom(userID uint, uuid string) (model.RoomInviteResponse, error)
	DeleteRoom(uuid string) error
	LeaveRoom(userID uint, roomUUID string) error
	UpdateMessage(messageUUID string, content string) (model.BroadcastMessage, error)
	DeleteMessage(messageUUID string) (model.BroadcastMessage, error)
	AddRoomChannel(roomUUID string) chan *pb.MessageResponse
	SendMessageToRoomChannel(roomUUID string, msg model.BroadcastMessage)
	DeleteRoomChannel(roomUUID string, ch chan *pb.MessageResponse)
	GetMessages(uuid string, offset uint) ([]*pb.MessageInfo, error)
}

type RoomUsecase struct {
	rr  repository.RoomRepositoryInterface
	rmr repository.RoomMemberRepositoryInterface
	ur  repository.UserRepositoryInterface
	fr  repository.FriendRepositoryInterface
	mr  repository.MessageRepositoryInterface
	db  *gorm.DB

	pb.UnimplementedMessageServiceServer
	msgChannels map[string]*model.RoomChannels
}

func NewRoomUsecase(
	rr repository.RoomRepositoryInterface,
	rmr repository.RoomMemberRepositoryInterface,
	ur repository.UserRepositoryInterface,
	fr repository.FriendRepositoryInterface,
	mr repository.MessageRepositoryInterface,
	db *gorm.DB,
) RoomUsecaseInterface {
	return &RoomUsecase{rr: rr,
		rmr:         rmr,
		ur:          ur,
		fr:          fr,
		mr:          mr,
		db:          db,
		msgChannels: make(map[string]*model.RoomChannels),
	}
}

func (ru *RoomUsecase) CreateDMRoom(userID uint, receiverID uint, tx *gorm.DB) error {
	uuid := xid.New().String()

	room := model.Room{
		UUID:        uuid,
		Name:        "",
		Type:        1,
		AdminUserID: 0,
	}

	if err := ru.rr.Insert(&room, tx); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	members := []model.RoomMember{
		{
			RoomID: room.ID,
			UserID: userID,
		},
		{
			RoomID: room.ID,
			UserID: receiverID,
		},
	}

	if err := ru.rmr.Insert(members, tx); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (ru RoomUsecase) CreateRoom(req model.RoomCreateRequest, userID uint) (model.RoomCreateResponse, error) {
	uuid := xid.New().String()

	room := model.Room{
		UUID:        uuid,
		Name:        req.Name,
		Type:        2,
		AdminUserID: userID,
	}
	tx := ru.db.Begin()
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return model.RoomCreateResponse{}, tx.Error
	}

	// トランザクション開始処理を実装する
	if err := ru.rr.Insert(&room, tx); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return model.RoomCreateResponse{}, err
	}

	members := []model.RoomMember{
		{
			RoomID: room.ID,
			UserID: userID,
		},
	}

	if err := ru.rmr.Insert(members, tx); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return model.RoomCreateResponse{}, err
	}

	// トランザクション終了処理を実装する
	if err := tx.Commit().Error; err != nil {
		fmt.Println(err)
		return model.RoomCreateResponse{}, err
	}

	return model.RoomCreateResponse{
		UUID: room.UUID,
		Name: room.Name,
	}, nil
}

func (ru RoomUsecase) GetRoomMessages(uuid string, userID uint, offset uint) (model.RoomInfoResponse, error) {
	room := model.Room{
		UUID: uuid,
	}

	// ルームレコードを取得
	if err := ru.rr.GetByUUID(&room); err != nil {
		fmt.Println(err)
		return model.RoomInfoResponse{}, err
	}

	// ルームメンバーレコードを取得
	roomMemberNames, err := ru.rmr.GetRoomMemberNamesByRoomID(room.ID)
	if err != nil {
		fmt.Println(err)
		return model.RoomInfoResponse{}, err
	}

	messages, err := ru.mr.GetMessagesByRoomID(room.ID, offset)
	if err != nil {
		fmt.Println(err)
		return model.RoomInfoResponse{}, err
	}

	isAdmin := false
	if room.AdminUserID == userID {
		isAdmin = true
	}

	sort.Slice(messages, func(i int, j int) bool {
		if messages[i].Timestamp.Equal(messages[j].Timestamp) {
			return messages[i].ID < messages[j].ID
		}
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	res := model.RoomInfoResponse{
		Name:     room.Name,
		UUID:     room.UUID,
		IsAdmin:  isAdmin,
		Members:  roomMemberNames,
		Messages: messages,
	}

	return res, nil
}

func (ru RoomUsecase) CreateMessage(roomUUID string, req model.MessageCreateRequest, userID uint) (model.BroadcastMessage, error) {
	room := model.Room{
		UUID: roomUUID,
	}
	if err := ru.rr.GetByUUID(&room); err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}

	uuid := xid.New().String()
	message := model.Message{
		UUID:    uuid,
		UserID:  userID,
		RoomID:  room.ID,
		Content: req.Content,
	}

	if err := ru.mr.Insert(&message); err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}

	if err := ru.mr.GetByID(&message); err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}

	room.LastMessageAt = message.CreatedAt
	if err := ru.rr.UpdateLastMessageAtByRoomUUID(&room); err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}

	user, err := ru.ur.GetUserNameByID(message.UserID)
	if err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}

	msg := model.BroadcastMessage{
		Type: "send",
		MessageInfo: model.MessageInfo{
			ID:        message.ID,
			UUID:      message.UUID,
			Content:   message.Content,
			Timestamp: message.CreatedAt,
			User: model.UserInfo{
				Name: user,
			},
		},
	}
	return msg, nil
}

func (ru RoomUsecase) GetRooms(userID uint) ([]model.GetRoomsResponse, error) {
	res, err := ru.rr.GetUUIDAndNameByRoomMemberUserID(userID, 2)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, nil
}

func (ru RoomUsecase) InviteRoom(userID uint, uuid string) (model.RoomInviteResponse, error) {
	room := &model.Room{
		UUID: uuid,
	}
	if err := ru.rr.GetByUUID(room); err != nil {
		fmt.Println(err)
		return model.RoomInviteResponse{}, err
	}
	member := []model.RoomMember{
		{
			UserID: userID,
			RoomID: room.ID,
		},
	}
	if err := ru.rmr.Insert(member, nil); err != nil {
		fmt.Println(err)
		return model.RoomInviteResponse{}, err
	}
	return model.RoomInviteResponse{UUID: room.UUID}, nil
}

func (ru RoomUsecase) DeleteRoom(uuid string) error {
	room := &model.Room{
		UUID: uuid,
	}
	if err := ru.rr.DeleteByRoomUUID(room); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (ru RoomUsecase) LeaveRoom(userID uint, roomUUID string) error {
	room := model.Room{
		UUID: roomUUID,
	}
	if err := ru.rr.GetByUUID(&room); err != nil {
		fmt.Println(err)
		return err
	}
	member := model.RoomMember{
		RoomID: room.ID,
		UserID: userID,
	}
	fmt.Println("LeaveRoom:", member)
	if err := ru.rmr.DeleteByRoomIDAndUserID(&member); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (ru RoomUsecase) UpdateMessage(messageUUID string, content string) (model.BroadcastMessage, error) {
	message := model.Message{
		UUID:    messageUUID,
		Content: content,
	}
	// メッセージ取得処理を実装
	if err := ru.mr.GetByUUID(&message); err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}
	message.Content = content
	if err := ru.mr.UpdateContentByUUID(&message); err != nil {
		fmt.Println(err)
		fmt.Println("ru.mr.UpdateMessageByMessageUUID: ", err)
		return model.BroadcastMessage{}, err
	}

	fmt.Println("UpdateMessage: ", message)
	mInfo, err := ru.mr.GetMessageByID(message.ID)
	if err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}
	msg := model.BroadcastMessage{
		Type:        "update",
		MessageInfo: mInfo,
	}
	return msg, nil
}

func (ru RoomUsecase) DeleteMessage(messageUUID string) (model.BroadcastMessage, error) {
	if err := ru.mr.DeleteByUUID(messageUUID); err != nil {
		fmt.Println(err)
		return model.BroadcastMessage{}, err
	}
	msg := model.BroadcastMessage{
		Type: "delete",
		MessageInfo: model.MessageInfo{
			UUID: messageUUID,
		},
	}
	return msg, nil
}

func (gu *RoomUsecase) AddRoomChannel(roomUUID string) chan *pb.MessageResponse {
	ch := make(chan *pb.MessageResponse)

	room, exists := gu.msgChannels[roomUUID]
	if !exists {
		room = &model.RoomChannels{
			ClientChannels: make([]chan *pb.MessageResponse, 0),
		}
		room.ClientChannels = append(room.ClientChannels, ch)
		gu.msgChannels[roomUUID] = room
		return ch
	}

	room.ClientChannels = append(room.ClientChannels, ch)
	return ch
}

func (gu *RoomUsecase) SendMessageToRoomChannel(roomUUID string, msg model.BroadcastMessage) {
	room, exists := gu.msgChannels[roomUUID]
	if !exists || room == nil {
		fmt.Println("Room does not exist or is nil")
		return
	}
	for i := range room.ClientChannels {
		ch := &room.ClientChannels[i]
		if *ch == nil {
			fmt.Println("Skipped nil channel")
			continue
		}
		res := &pb.MessageResponse{
			Type: msg.Type,
			MessageInfo: &pb.MessageInfo{
				Id:        uint32(msg.MessageInfo.ID),
				Uuid:      msg.MessageInfo.UUID,
				Content:   msg.MessageInfo.Content,
				Timestamp: msg.MessageInfo.Timestamp.String(),
				User: &pb.UserInfo{
					Name: msg.MessageInfo.User.Name,
				},
			},
		}
		*ch <- res
		fmt.Println("SendMessageToRoomChannel ch: ", ch)
	}
}

func (gu *RoomUsecase) DeleteRoomChannel(roomUUID string, ch chan *pb.MessageResponse) {
	room := gu.msgChannels[roomUUID]
	room.Mu.Lock()
	defer room.Mu.Unlock()
	for i, c := range room.ClientChannels {
		if c == ch {
			room.ClientChannels = append(room.ClientChannels[:i], room.ClientChannels[i+1:]...)
			break
		}
	}

	if len(room.ClientChannels) == 0 {
		delete(gu.msgChannels, roomUUID)
		fmt.Println("Room removed from gu.msgChannels:", gu.msgChannels)
	}

	close(ch)
}

func (ru RoomUsecase) GetMessages(uuid string, offset uint) ([]*pb.MessageInfo, error) {
	room := model.Room{
		UUID: uuid,
	}

	// ルームレコードを取得
	if err := ru.rr.GetByUUID(&room); err != nil {
		fmt.Println(err)
		return nil, err
	}

	messages, err := ru.mr.GetMessagesByRoomID(room.ID, offset)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sort.Slice(messages, func(i int, j int) bool {
		if messages[i].Timestamp.Equal(messages[j].Timestamp) {
			return messages[i].ID < messages[j].ID
		}
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	var res []*pb.MessageInfo
	for _, m := range messages {
		pbm := &pb.MessageInfo{
			Id:        uint32(m.ID),
			Uuid:      m.UUID,
			Content:   m.Content,
			Timestamp: m.Timestamp.String(),
			User: &pb.UserInfo{
				Name: m.User.Name,
			},
		}
		res = append(res, pbm)
	}

	return res, nil
}
