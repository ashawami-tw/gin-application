package model

import "github.com/google/uuid"

type UserName struct {
	ID        uuid.UUID `gorm:"primarykey;default:uuid_generate_v4()" sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserId    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func (*UserName) TableName() string {
	return "username"
}
