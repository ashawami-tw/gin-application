package username

import (
	"github.com/jackc/pgconn"
	"net/http"
	"some-application/backend/model"
	"some-application/backend/repository"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
)

type Service interface {
	AddUserName(username *model.UserName) *utils.Error
}

type service struct {
	usernameRepo repository.UserNameRepository
}

func NewUsernameService(usernameRepo repository.UserNameRepository) Service {
	return &service{
		usernameRepo: usernameRepo,
	}
}

func (s *service) AddUserName(username *model.UserName) *utils.Error {
	err := s.usernameRepo.InsertName(username)
	if err != nil && err.(*pgconn.PgError).ConstraintName == "unique_user_id" {
		return utils.LogError(http.StatusBadRequest, err.Error(), constant.NameAlreadyAdded)
	}
	if err != nil {
		return utils.LogError(http.StatusInternalServerError, err.Error(), constant.InternalServerError)
	}
	return nil
}
