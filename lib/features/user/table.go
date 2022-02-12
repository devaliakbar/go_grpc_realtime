package user

import "go_grpc_realtime/lib/core/database"

func migrateDb() {
	database.DB.AutoMigrate(&User{})
}

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
}
