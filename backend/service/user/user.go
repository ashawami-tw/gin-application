package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"some-application/backend/model"
	"some-application/backend/repository"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
)

type Service interface {
	AddUser(user *model.User) *utils.Error
	HashPassword(user *model.User) *utils.Error
	UserAlreadyExists(newEmail string) (*utils.Error, bool)
	UserIdExists(id uuid.UUID) (*utils.Error, bool)
}

type service struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) AddUser(user *model.User) *utils.Error {
	err := s.userRepo.InsertEmail(user)
	if err != nil {
		return utils.LogError(http.StatusInternalServerError, err.Error(), constant.InternalServerError)
	}
	return nil
}

func (s *service) HashPassword(user *model.User) *utils.Error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), constant.PasswordCost)
	if err != nil {
		return utils.LogError(http.StatusInternalServerError, err.Error(), constant.InternalServerError)
	}
	user.Password = string(hashPassword[:])
	return nil
}

func (s *service) UserAlreadyExists(newEmail string) (*utils.Error, bool) {
	user, err := s.userRepo.FindUserByEmail(newEmail)
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.LogError(http.StatusInternalServerError, err.Error(), constant.InternalServerError), true
	}
	if user.Id != uuid.Nil {
		return utils.LogError(http.StatusBadRequest, "user already exists", constant.EmailAlreadyExists), true
	}
	return nil, false
}

func (s *service) UserIdExists(id uuid.UUID) (*utils.Error, bool) {
	_, err := s.userRepo.FindUserById(id)
	if err == gorm.ErrRecordNotFound {
		return utils.LogError(http.StatusBadRequest, err.Error(), constant.InvalidUserId), false
	}
	if err != nil {
		return utils.LogError(http.StatusInternalServerError, err.Error(), constant.InternalServerError), false
	}
	return nil, true
}
