package message

import (
	"context"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/interceptors"
	"go_grpc_realtime/lib/core/jwtmanager"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessageController struct {
	mu sync.Mutex
	grpcgen.UnimplementedMessageServiceServer
	*repository
	messageListeners map[uuid.UUID]*MessageListener
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func InitAndGetMessageServices() grpcgen.MessageServiceServer {
	repo := &repository{}

	repo.migrateDb()

	return &MessageController{
		repository:       repo,
		messageListeners: map[uuid.UUID]*MessageListener{},
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ctr *MessageController) GetMessageRooms(ctx context.Context, req *grpcgen.GetMessageRoomsRequest) (*grpcgen.GetMessageRoomsResponse, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"user not found",
		)
	}

	return ctr.getMessageRooms(req, userId)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ctr *MessageController) GetMessageRoomDetails(ctx context.Context, req *grpcgen.GetMessageRoomDetailsRequest) (*grpcgen.MessageRoom, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"user not found",
		)
	}

	return ctr.getMessageRoomDetails(req, userId)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ctr *MessageController) GetMessages(ctx context.Context, req *grpcgen.GetMessagesRequest) (*grpcgen.GetMessagesResponse, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"user not found",
		)
	}

	return ctr.getMessages(req, userId)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ctr *MessageController) SendMessage(ctx context.Context, req *grpcgen.SendMessageRequest) (*grpcgen.Message, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"user not found",
		)
	}

	roomMembers, message, err := ctr.sendMessage(req, userId)

	for _, mem := range roomMembers {
		ctr.mu.Lock()
		for _, msgListener := range ctr.messageListeners {
			if msgListener.UserId == mem.UserId {
				msgListener.Channel <- message
			}
		}
		ctr.mu.Unlock()
	}

	return message, err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ctr *MessageController) ListenToNewMessage(req *grpcgen.ListenToNewMessageRequest, stream grpcgen.MessageService_ListenToNewMessageServer) error {
	uid, err := interceptors.GetUserIdFromHeader(stream.Context())
	if err != nil {
		return err
	}

	ctr.mu.Lock()
	key := uuid.New()
	msgChannel := make(chan *grpcgen.Message)
	ctr.messageListeners[key] = &MessageListener{
		UserId:  uid,
		Channel: msgChannel,
	}
	ctr.mu.Unlock()

	for {
		select {
		case <-stream.Context().Done():
			ctr.mu.Lock()
			delete(ctr.messageListeners, key)
			ctr.mu.Unlock()

			return nil
		case msg := <-msgChannel:
			stream.Send(&grpcgen.ListenToNewMessageResponse{
				NewMessage: msg,
			})
		}
	}
}
