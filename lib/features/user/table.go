package user

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	FullName string `json:"full_name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
}
