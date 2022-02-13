package message

import (
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/utils"
)

type repository struct {
	*utils.Validation
}

func (*repository) migrateDb() {
	database.DB.AutoMigrate(&RoomTbl{}, &RoomMembersTbl{}, &MessageTbl{})
}
