package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"some-application/backend/model"
	"strings"
)

type UserRepository interface {
	InsertEmail(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserById(id uuid.UUID) (*model.User, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (d *gormUserRepository) InsertEmail(user *model.User) error {
	user.Email = strings.ToLower(user.Email)
	return d.db.Table(user.TableName()).Model(&model.User{}).Omit("id").Create(&user).Error
}

func (d *gormUserRepository) FindUserByEmail(email string) (*model.User, error) {
	var existingUser model.User
	email = strings.ToLower(email)
	if err := d.db.Table(existingUser.TableName()).Where("email = ?", email).First(&existingUser).Error; err != nil {
		return &existingUser, err
	}
	return &existingUser, nil
}

func (d *gormUserRepository) FindUserById(id uuid.UUID) (*model.User, error) {
	var existingUser model.User
	if err := d.db.Table(existingUser.TableName()).Where("id = ?", id).First(&existingUser).Error; err != nil {
		return &existingUser, err
	}
	return &existingUser, nil
}
