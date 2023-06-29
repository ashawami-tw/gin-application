package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"some-application/backend/model"
	"some-application/backend/repository"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
	"testing"
)

func setupTest() (*repository.MockUserRepository, Service) {
	mockUserRepo := repository.NewUserRepositoryMock()
	userService := NewUserService(mockUserRepo)

	return mockUserRepo, userService
}

func TestService_UserAlreadyExists(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.New(), "temp@gmail.com", "password")
	respErr := utils.LogError(http.StatusBadRequest, "user already exists", constant.EmailAlreadyExists)
	mockUserRepo.Mock.On("FindUserByEmail", user.Email).Return(user, gorm.ErrRecordNotFound)
	err, emailExists := userService.UserAlreadyExists(user.Email)

	assert.Equalf(t, err, respErr, "Error should be retruned as user already exists")
	assert.Equalf(t, true, emailExists, "service should return true as user already exists")
}

func TestService_UserDoesNotExists(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.Nil, "temp@gmail.com", "password")

	mockUserRepo.Mock.On("FindUserByEmail", user.Email).Return(user, nil)
	err, emailExists := userService.UserAlreadyExists(user.Email)

	assert.Nil(t, err, "Error should be nil")
	assert.Equalf(t, false, emailExists, "service should return false as user does not exists")
}

func TestService_UserAlreadyExists_InternalServerError(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.Nil, "", "password")
	dbError := errors.New("error while connecting to DB")
	respErr := utils.LogError(http.StatusInternalServerError, dbError.Error(), constant.InternalServerError)

	mockUserRepo.Mock.On("FindUserByEmail", "").Return(user, dbError)
	err, emailExists := userService.UserAlreadyExists("")

	assert.Equalf(t, respErr, err, "Error should be returned as connection to DB fails")
	assert.Equalf(t, true, emailExists, "service should return false as error occured while connecting to DB")
}

func TestService_HashPassword(t *testing.T) {
	user := getUser(uuid.Nil, "", "password")
	_, userService := setupTest()
	err := userService.HashPassword(user)
	assert.Nil(t, err)
}

func TestService_AddEmail(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.Nil, "temp@gmail.com", "password")
	mockUserRepo.Mock.On("InsertEmail", user).Return(nil)
	err := userService.AddUser(user)

	assert.Nil(t, err, "Error should be nil")
}

func TestService_AddUser_InternalServerError(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.Nil, "temp@gmail.com", "password")
	dbError := errors.New("error while connecting to DB")
	respErr := utils.LogError(http.StatusInternalServerError, dbError.Error(), constant.InternalServerError)
	mockUserRepo.Mock.On("InsertEmail", user).Return(dbError)

	err := userService.AddUser(user)

	assert.Equalf(t, respErr, err, "Error should be returned as connection to DB fails")
}

func TestService_UserIdExists(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.New(), "temp@gmail.com", "password")

	mockUserRepo.Mock.On("FindUserById", user.Id).Return(user, nil)
	err, idExists := userService.UserIdExists(user.Id)

	assert.Nil(t, err, "Error should be nil")
	assert.Equalf(t, true, idExists, "service should return true as user id exists")
}

func TestService_UserIdDoesNotExists(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.New(), "temp@gmail.com", "password")
	dbError := gorm.ErrRecordNotFound
	respErr := utils.LogError(http.StatusBadRequest, dbError.Error(), constant.InvalidUserId)

	mockUserRepo.Mock.On("FindUserById", user.Id).Return(user, dbError)
	err, idExists := userService.UserIdExists(user.Id)

	assert.Equalf(t, respErr, err, "Error should be returned when user id does not exists")
	assert.Equalf(t, idExists, false, "service should return false as user id does not exist")
}

func TestService_UserIdExists_InternalServerError(t *testing.T) {
	mockUserRepo, userService := setupTest()
	user := getUser(uuid.New(), "temp@gmail.com", "password")
	dbError := errors.New("error while connecting to DB")
	respError := utils.LogError(http.StatusInternalServerError, dbError.Error(), constant.InternalServerError)

	mockUserRepo.Mock.On("FindUserById", user.Id).Return(user, dbError)
	err, idExists := userService.UserIdExists(user.Id)

	assert.Equalf(t, respError, err, "Error should be returned as connection to DB fails")
	assert.Equalf(t, idExists, false, "service should return false as connection to DB fails")
}

func getUser(id uuid.UUID, email, password string) *model.User {
	return &model.User{
		Email:    email,
		Password: password,
		Id:       id,
	}
}
