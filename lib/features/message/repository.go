package message

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/utils"
	"go_grpc_realtime/lib/features/user"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repository struct {
	*utils.Validation
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (*repository) migrateDb() {
	database.DB.AutoMigrate(&RoomTbl{}, &RoomMembersTbl{}, &MessageTbl{})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (repo *repository) createMessageRoom(req *grpcgen.CreateMessageRoomRequest, uid uint) (*grpcgen.MessageRoom, error) {
	///Converting string member id to uid and removing duplicate userid
	membersArr := req.GetMembers()
	membersArr = append(membersArr, fmt.Sprint(uid))
	members, err := utils.RemoveDuplicateUsers(membersArr)
	if err != nil {
		return nil, err
	}

	if len(members) == 1 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"members are empty",
		)
	}
	if !req.IsOneToOne && len(members) < 3 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"atleast 2 members are required to create group",
		)
	}
	roomName := strings.TrimSpace(req.GetRoomName())
	if !req.IsOneToOne && roomName == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"group name is empty",
		)
	}

	///If one to one, then check room already exist or not, if exist simply gives back the already existed room
	if req.GetIsOneToOne() {
		alrExistRm := checkRoomAlreadyExist(uid, members)
		if alrExistRm != nil {
			return alrExistRm, nil
		}
	}

	///Create room
	if req.IsOneToOne {
		//This will help in future to check weather that particular one to one room already exist or not
		roomName = fmt.Sprintf("%d-%d", members[0], members[1])
	}

	transactionDb := database.DB.Begin()

	createdAt := time.Now()

	room := RoomTbl{
		Name:       roomName,
		IsOneToOne: req.GetIsOneToOne(),
	}
	if !req.IsOneToOne {
		room.LastUpdated = &createdAt
	}
	if crtRmErr := transactionDb.Create(&room).Error; crtRmErr != nil {
		transactionDb.Rollback()
		return nil, status.Errorf(
			codes.Internal,
			crtRmErr.Error(),
		)
	}

	var grpMem []*grpcgen.User
	///Checking all passed userid existed and if exit add to room members tables
	for _, mem := range members {
		var usr user.UserTbl
		if err := database.DB.Where("id = ?", mem).First(&usr).Error; err != nil {
			transactionDb.Rollback()
			return nil, status.Errorf(
				codes.NotFound,
				"User not found",
			)
		}

		///Adding members to room
		rmMem := RoomMembersTbl{
			RoomId: room.ID,
			UserId: mem,
		}
		if addMemErr := transactionDb.Create(&rmMem).Error; addMemErr != nil {
			transactionDb.Rollback()
			return nil, status.Errorf(
				codes.Internal,
				addMemErr.Error(),
			)
		}

		if mem != uid {
			grpMem = append(grpMem, &grpcgen.User{
				Id:       fmt.Sprint(usr.ID),
				FullName: usr.FullName,
				Email:    usr.Email,
			})
		}
	}

	///If not One-To-One, add a message regarding who created group
	if !req.IsOneToOne {
		grpMsg := MessageTbl{
			RoomId:    room.ID,
			SenderId:  uid,
			Body:      "created group",
			CreatedAt: createdAt,
		}
		if addMesErr := transactionDb.Create(&grpMsg).Error; addMesErr != nil {
			transactionDb.Rollback()
			return nil, status.Errorf(
				codes.Internal,
				addMesErr.Error(),
			)
		}
	}

	if req.IsOneToOne {
		roomName = grpMem[0].FullName
	}

	transactionDb.Commit()
	return &grpcgen.MessageRoom{
		Id:         fmt.Sprint(room.ID),
		RoomName:   roomName,
		IsOneToOne: room.IsOneToOne,
		Members:    grpMem,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkRoomAlreadyExist(uid uint, membs []uint) *grpcgen.MessageRoom {
	possibleName1 := fmt.Sprintf("%d-%d", membs[0], membs[1])
	possibleName2 := fmt.Sprintf("%d-%d", membs[1], membs[0])
	room := RoomTbl{}
	if err := database.DB.Where("(name = ? OR name = ?) AND is_one_to_one = true", possibleName1, possibleName2).First(&room).Error; err != nil {
		return nil
	}

	userId := membs[0]
	if uid == membs[0] {
		userId = membs[1]
	}

	var usr user.UserTbl
	if err := database.DB.Where("id = ?", userId).First(&usr).Error; err != nil {
		return nil
	}

	return &grpcgen.MessageRoom{
		Id:         fmt.Sprint(room.ID),
		RoomName:   usr.FullName,
		IsOneToOne: true,
		Members: []*grpcgen.User{
			{
				Id:       fmt.Sprint(usr.ID),
				FullName: usr.FullName,
				Email:    usr.Email,
			},
		},
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (repo *repository) getMessageRooms(req *grpcgen.GetMessageRoomsRequest, uid uint) (*grpcgen.GetMessageRoomsResponse, error) {
	skip := int(req.GetSkip())

	take := 10
	if req.GetTake() != 0 {
		take = int(req.GetTake())
		if take > 100 {
			take = 100
		}
	}

	roomsArr := []MessageRoomQuery{}

	database.DB.Table("room_tbls").
		Joins("inner join room_members_tbls on room_members_tbls.room_id = room_tbls.id").
		Select("room_tbls.id as id,room_tbls.name as room_name, room_tbls.is_one_to_one as is_one_to_one").
		Order("room_tbls.last_updated desc").
		Offset(skip).Limit(take).
		Find(&roomsArr, "room_tbls.last_updated IS NOT NULL and room_members_tbls.user_id = ?", uid)

	rooms := []*grpcgen.MessageRoom{}
	for _, room := range roomsArr {
		roomName := room.RoomName
		roomMembs := []*grpcgen.User{}

		if room.IsOneToOne || req.GetGetGroupMembers() {
			members := []user.UserQuery{}
			database.DB.Table("user_tbls").
				Joins("inner join room_members_tbls on room_members_tbls.user_id = user_tbls.id").
				Select("user_tbls.id as id,user_tbls.full_name as full_name, user_tbls.email as email").
				Order("user_tbls.full_name asc").
				Find(&members, "room_members_tbls.room_id = ?", room.ID)

			for _, member := range members {
				if member.ID != uid {
					///If chat is one to one, set room name to end user name
					if room.IsOneToOne {
						roomName = member.FullName
					}
					roomMembs = append(roomMembs, &grpcgen.User{
						Id:       fmt.Sprint(member.ID),
						FullName: member.FullName,
						Email:    member.Email,
					})
				}
			}
		}

		rooms = append(rooms, &grpcgen.MessageRoom{
			Id:         fmt.Sprint(room.ID),
			RoomName:   roomName,
			IsOneToOne: room.IsOneToOne,
			Members:    roomMembs,
		})
	}

	return &grpcgen.GetMessageRoomsResponse{
		Rooms: rooms,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (repo *repository) getMessageRoomDetails(req *grpcgen.GetMessageRoomDetailsRequest, uid uint) (*grpcgen.MessageRoom, error) {
	members := []user.UserQuery{}
	database.DB.Table("user_tbls").
		Joins("inner join room_members_tbls on room_members_tbls.user_id = user_tbls.id").
		Select("user_tbls.id as id,user_tbls.full_name as full_name, user_tbls.email as email").
		Order("user_tbls.full_name asc").
		Find(&members, "room_members_tbls.room_id = ?", req.GetRoomId())

	if len(members) < 2 {
		return nil, status.Errorf(
			codes.NotFound,
			"room not found",
		)
	}

	roomMembs := []*grpcgen.User{}

	isCurrentUserExistInThisRoom := false
	for _, member := range members {
		if member.ID == uid {
			isCurrentUserExistInThisRoom = true
		} else {
			roomMembs = append(roomMembs, &grpcgen.User{
				Id:       fmt.Sprint(member.ID),
				FullName: member.FullName,
				Email:    member.Email,
			})
		}
	}

	if !isCurrentUserExistInThisRoom {
		return nil, status.Errorf(
			codes.NotFound,
			"room not found",
		)
	}

	var room RoomTbl
	if err := database.DB.Where("id = ?", req.GetRoomId()).First(&room).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"room not found",
		)
	}

	roomName := room.Name
	if room.IsOneToOne {
		roomName = roomMembs[0].FullName
	}

	return &grpcgen.MessageRoom{
		Id:         fmt.Sprint(room.ID),
		RoomName:   roomName,
		IsOneToOne: room.IsOneToOne,
		Members:    roomMembs,
	}, nil
}
