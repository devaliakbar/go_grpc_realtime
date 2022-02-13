package message

import (
	"context"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/jwtmanager"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessageController struct {
	grpcgen.UnimplementedMessageServiceServer
	*repository
}

func InitAndGetMessageServices() grpcgen.MessageServiceServer {
	repo := &repository{}

	repo.migrateDb()

	return &MessageController{
		repository: repo,
	}
}

func (ctr *MessageController) CreateMessageRoom(ctx context.Context, req *grpcgen.CreateMessageRoomRequest) (*grpcgen.MessageRoom, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"user not found",
		)
	}

	return ctr.createMessageRoom(req, userId)
}

func (ctr *MessageController) GetMessageRooms(ctx context.Context, req *grpcgen.GetMessageRoomsRequest) (*grpcgen.GetMessageRoomsResponse, error) {
	log.Println("GetRooms")
	return nil, nil
}

func (ctr *MessageController) GetMessageRoomMembers(ctx context.Context, req *grpcgen.GetMessageRoomMembersRequest) (*grpcgen.GetMessageRoomMembersResponse, error) {
	log.Println("GetRoomsMem")
	return nil, nil
}

func (ctr *MessageController) SendMessage(ctx context.Context, req *grpcgen.SendMessageRequest) (*grpcgen.Message, error) {
	log.Println("SendMes")
	return nil, nil
}

func (ctr *MessageController) GetMessages(ctx context.Context, req *grpcgen.GetMessagesRequest) (*grpcgen.GetMessagesResponse, error) {
	log.Println("Get mes")
	return nil, nil
}
