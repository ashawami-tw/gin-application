package username

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"net/http"
	"some-application/backend/model"
	"some-application/backend/repository"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"

	"testing"
)

func setupTest() (*repository.MockUsernameRepository, Service) {
	mockUsernameRepo := repository.NewUsernameRepositoryMock()
	usernameService := NewUsernameService(mockUsernameRepo)
	return mockUsernameRepo, usernameService
}

func TestService_AddUserName(t *testing.T) {
	mockUsernameRepo, usernameService := setupTest()
	username := getUsername(uuid.Nil, uuid.New(), "Steve", "Smith")

	mockUsernameRepo.Mock.On("InsertName", username).Return(nil)

	err := usernameService.AddUserName(username)

	assert.Nil(t, err, "Error should be nil")
}

func TestService_AddUserName_InternalServerError(t *testing.T) {
	mockUsernameRepo, usernameService := setupTest()
	username := getUsername(uuid.Nil, uuid.New(), "Steve", "Smith")
	pgError := &pgconn.PgError{Severity: "ERROR", Message: "error while connecting to DB", Code: "23505"}
	dbError := errors.New("ERROR: error while connecting to DB (SQLSTATE 23505)")
	resError := utils.LogError(http.StatusInternalServerError, dbError.Error(), constant.InternalServerError)

	mockUsernameRepo.Mock.On("InsertName", username).Return(pgError)

	err := usernameService.AddUserName(username)

	assert.Equalf(t, resError, err, "Error should be returned as connection to DB fails")
}

func TestService_AddUserName_UniqueUserId(t *testing.T) {
	mockUsernameRepo, usernameService := setupTest()
	username := getUsername(uuid.Nil, uuid.New(), "Steve", "Smith")
	pgError := &pgconn.PgError{ConstraintName: "unique_user_id", Severity: "ERROR", Message: "duplicate key value violates unique constraint unique_user_id", Code: "23505"}
	dbError := errors.New("ERROR: duplicate key value violates unique constraint unique_user_id (SQLSTATE 23505)")
	resError := utils.LogError(http.StatusBadRequest, dbError.Error(), constant.NameAlreadyAdded)

	mockUsernameRepo.Mock.On("InsertName", username).Return(pgError)

	err := usernameService.AddUserName(username)
	assert.Equalf(t, resError, err, "Error should be returned when first name and last name are already added for the user")
}

func getUsername(id uuid.UUID, userId uuid.UUID, firstName, lastName string) *model.UserName {
	return &model.UserName{
		ID:        id,
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
	}
}
