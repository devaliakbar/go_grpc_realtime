package message

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/utils"
	"go_grpc_realtime/lib/features/user"
	"strings"

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

	room := RoomTbl{
		Name:       roomName,
		IsOneToOne: req.GetIsOneToOne(),
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
