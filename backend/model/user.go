package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `gorm:"primarykey;default:uuid_generate_v4()" sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (*User) TableName() string {
	return "user"
}
