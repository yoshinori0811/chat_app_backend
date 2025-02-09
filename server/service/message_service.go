package service

import (
	"context"
	"fmt"

	pb "github.com/yoshinori0811/chat_app_backend/pb"
	"github.com/yoshinori0811/chat_app_backend/usecase"
)

type MessageServiceServer struct {
	pb.UnimplementedMessageServiceServer
	ru usecase.RoomUsecaseInterface
}

func NewMessageServiceServer(ru usecase.RoomUsecaseInterface) *MessageServiceServer {
	return &MessageServiceServer{
		ru: ru,
	}
}

func (m *MessageServiceServer) Connect(req *pb.ConnectRequest, stream pb.MessageService_ConnectServer) error {
	ctx := stream.Context()

	uuid := req.Uuid
	ch := m.ru.AddRoomChannel(uuid)

	defer m.ru.DeleteRoomChannel(uuid, ch)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Close ch:", ch)
			return ctx.Err()

		case msg := <-ch:
			if err := stream.Send(msg); err != nil {
				fmt.Println("Error stream message:", err)
				return err
			}
		}
	}
}

func (m *MessageServiceServer) GetMessages(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessagesResponse, error) {
	uuid := req.Uuid
	offset := req.Offset

	res, err := m.ru.GetMessages(uuid, uint(offset))
	if err != nil {
		return nil, err
	}

	return &pb.GetMessagesResponse{
		Messages: res,
	}, nil
}
