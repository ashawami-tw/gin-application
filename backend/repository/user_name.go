package repository

import (
	"gorm.io/gorm"
	"some-application/backend/model"
	"strings"
)

type UserNameRepository interface {
	InsertName(username *model.UserName) error
}

type gormUserNameRepository struct {
	db *gorm.DB
}

func NewUserNameRepository(db *gorm.DB) UserNameRepository {
	return &gormUserNameRepository{
		db: db,
	}
}

func (d *gormUserNameRepository) InsertName(username *model.UserName) error {
	username.FirstName = strings.ToLower(username.FirstName)
	username.LastName = strings.ToLower(username.LastName)
	return d.db.Table(username.TableName()).Model(&model.UserName{}).Omit("id").Create(&username).Error
}
