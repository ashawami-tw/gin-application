package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"some-application/backend/kafka"
	"some-application/backend/kafka/message"
	"some-application/backend/model"
	"some-application/backend/service/user"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
	"testing"
)

func TestUserHandler_Success(t *testing.T) {
	type args struct {
		err                error
		expectedStatusCode int
		message            string
		userExists         bool
		mockUserResp       *model.User
	}

	successfulUserCreation := args{
		err:                errors.New("error while producing kafka event"),
		expectedStatusCode: http.StatusCreated,
		message:            constant.UserAddedSuccessfully,
		userExists:         false,
		mockUserResp: &model.User{
			Email:    "success@gmail.com",
			Password: "password",
		},
	}

	tests := []struct {
		name         string
		args         args
		mockProducer *kafka.MockProducer
	}{
		{
			name: "error while producing kafka error",
			args: successfulUserCreation,
			mockProducer: func() *kafka.MockProducer {
				mockProducer := kafka.NewMockProducer()
				mockProducer.Mock.On("EmitEvent", kafka.NewUserEvent, kafka.UserTopic, successfulUserCreation.mockUserResp.Id.String(), message.NewUser{
					Email: successfulUserCreation.mockUserResp.Email,
				}).Return(successfulUserCreation.err)
				return mockProducer
			}(),
		},
		{
			name: "success",
			args: successfulUserCreation,
			mockProducer: func() *kafka.MockProducer {
				mockProducer := kafka.NewMockProducer()
				mockProducer.Mock.On("EmitEvent", kafka.NewUserEvent, kafka.UserTopic, successfulUserCreation.mockUserResp.Id.String(), message.NewUser{
					Email: successfulUserCreation.mockUserResp.Email,
				}).Return(nil)
				return mockProducer
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// dependency
			mockUserService := user.NewUserServiceMock()
			userHandler := NewUserHandler(mockUserService, test.mockProducer)

			// prepare user req
			jsonUserResp, _ := json.Marshal(test.args.mockUserResp)
			req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonUserResp))

			// service mocks
			mockUserService.Mock.On("UserAlreadyExists", test.args.mockUserResp.Email).Return(nil, test.args.userExists)
			mockUserService.Mock.On("HashPassword", test.args.mockUserResp).Return(nil)
			mockUserService.Mock.On("AddUser", test.args.mockUserResp).Return(nil)

			// setup to call handler
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = req
			userHandler.CreateUser(ctx)

			// response
			var actualResponse utils.Response
			_ = json.NewDecoder(rec.Body).Decode(&actualResponse)

			// assertion
			assert.Equal(t, test.args.expectedStatusCode, rec.Code)
			assert.Equal(t, test.args.expectedStatusCode, actualResponse.Status)
			assert.Equal(t, test.args.message, actualResponse.Message)
		})
	}
}

func TestNewUserHandler_Failure(t *testing.T) {
	type args struct {
		err                error
		expectedStatusCode int
		message            string
		userExists         bool
		mockUserResp       *model.User
	}

	userAlreadyExists := args{
		err:                errors.New("user already exists"),
		expectedStatusCode: http.StatusBadRequest,
		message:            constant.EmailAlreadyExists,
		userExists:         true,
		mockUserResp: &model.User{
			Email:    "userExists@gmail.com",
			Password: "password",
		},
	}

	dbError := args{
		err:                errors.New("error while connecting to DB"),
		expectedStatusCode: http.StatusInternalServerError,
		message:            constant.InternalServerError,
		userExists:         true,
		mockUserResp: &model.User{
			Email:    "DBError@gmail.com",
			Password: "password",
		},
	}

	hashPasswordError := args{
		err:                errors.New("error while hashing password"),
		expectedStatusCode: http.StatusInternalServerError,
		message:            constant.InternalServerError,
		userExists:         false,
		mockUserResp: &model.User{
			Email:    "hashPassword@gmail.com",
			Password: "password",
		},
	}

	tests := []struct {
		name        string
		args        args
		mockService *user.MockService
	}{
		{
			name: "user already exists",
			args: userAlreadyExists,
			mockService: func() *user.MockService {
				mockUserService := user.NewUserServiceMock()
				err := utils.LogError(userAlreadyExists.expectedStatusCode, userAlreadyExists.err.Error(), userAlreadyExists.message)
				mockUserService.Mock.On("UserAlreadyExists", userAlreadyExists.mockUserResp.Email).Return(err, userAlreadyExists.userExists)
				return mockUserService
			}(),
		},
		{
			name: "error while querying DB to check if user exists",
			args: dbError,
			mockService: func() *user.MockService {
				mockUserService := user.NewUserServiceMock()
				err := utils.LogError(dbError.expectedStatusCode, dbError.err.Error(), dbError.message)
				mockUserService.Mock.On("UserAlreadyExists", dbError.mockUserResp.Email).Return(err, dbError.userExists)
				return mockUserService
			}(),
		},
		{
			name: "error while hashing password",
			args: hashPasswordError,
			mockService: func() *user.MockService {
				mockUserService := user.NewUserServiceMock()
				err := utils.LogError(hashPasswordError.expectedStatusCode, hashPasswordError.err.Error(), hashPasswordError.message)
				mockUserService.Mock.On("UserAlreadyExists", hashPasswordError.mockUserResp.Email).Return(nil, hashPasswordError.userExists)
				mockUserService.Mock.On("HashPassword", hashPasswordError.mockUserResp).Return(err)
				return mockUserService
			}(),
		},
		{
			name: "error while inserting email in DB",
			args: dbError,
			mockService: func() *user.MockService {
				mockUserService := user.NewUserServiceMock()
				dbError.userExists = false
				err := utils.LogError(dbError.expectedStatusCode, dbError.err.Error(), dbError.message)
				mockUserService.Mock.On("UserAlreadyExists", dbError.mockUserResp.Email).Return(nil, dbError.userExists)
				mockUserService.Mock.On("HashPassword", dbError.mockUserResp).Return(nil)
				mockUserService.Mock.On("AddUser", dbError.mockUserResp).Return(err)
				return mockUserService
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// dependency
			mockProducer := kafka.NewMockProducer()
			userHandler := NewUserHandler(test.mockService, mockProducer)

			// prepare user req
			jsonBody, _ := json.Marshal(test.args.mockUserResp)
			req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonBody))

			// setup to call handler
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = req
			userHandler.CreateUser(ctx)

			// response
			var actualResponse utils.Response
			_ = json.NewDecoder(rec.Body).Decode(&actualResponse)

			// assertion
			assert.Equal(t, test.args.expectedStatusCode, rec.Code)
			assert.Equal(t, test.args.expectedStatusCode, actualResponse.Status)
			assert.Equal(t, test.args.message, actualResponse.Message)
		})
	}
}

func TestUserHandler_Validation(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		message            string
		mockUserReq        *model.User
	}{
		{
			name:               "Empty request body",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUserReq:        &model.User{},
		},
		{
			name:               "Empty email",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUserReq: &model.User{
				Email:    "",
				Password: "password",
			},
		},
		{
			name:               "Empty password",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUserReq: &model.User{
				Email:    "test@gmail.com",
				Password: "",
			},
		},
		{
			name:               "Email without domain name",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUserReq: &model.User{
				Email:    "test.com",
				Password: "password",
			},
		},
		{
			name:               "Email without extension",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUserReq: &model.User{
				Email:    "test@gmail",
				Password: "password",
			},
		},
		{
			name:               "Email without username",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUserReq: &model.User{
				Email:    "@gmail.com",
				Password: "password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// dependency
			mockService := user.NewUserServiceMock()
			mockProducer := kafka.NewMockProducer()
			userHandler := NewUserHandler(mockService, mockProducer)

			// prepare user req
			jsonBody, _ := json.Marshal(test.mockUserReq)
			req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonBody))

			// setup to call handler
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = req
			userHandler.CreateUser(ctx)

			// response
			var actualResp utils.Response
			_ = json.NewDecoder(rec.Body).Decode(&actualResp)

			// assertion
			assert.Equal(t, test.expectedStatusCode, rec.Code)
			assert.Equal(t, test.expectedStatusCode, actualResp.Status)
			assert.Equal(t, test.message, actualResp.Message)

		})
	}
}

// other approach to write tests -

//func TestCreateUser_Failure(t *testing.T) {
//	route := gin.Default()
//	mockUserService := user.NewUserServiceMock()
//	mockProducer := kafka.NewMockProducer()
//	userHandler := NewUserHandler(mockUserService, mockProducer)
//	route.POST("/create", userHandler.CreateUser)
//	url := "/create"
//	mockUserReq := &model.User{
//		Email:    "test@gmail.com",
//		Password: "test",
//	}
//	var actualResponse utils.Response
//
//	t.Run("user already exists", func(t *testing.T) {
//		err := utils.LogError(http.StatusBadRequest, "user already exists", constant.EmailAlreadyExists)
//		mockUserService.Mock.On("UserAlreadyExists", mockUserReq.Email).Return(err, true).Once()
//
//		jsonUserResp, _ := json.Marshal(mockUserReq)
//		mockUserReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonUserResp))
//		rec := httptest.NewRecorder()
//		route.ServeHTTP(rec, mockUserReq)
//
//		_ = json.NewDecoder(rec.Body).Decode(&actualResponse)
//
//		assert.Equal(t, http.StatusBadRequest, rec.Code)
//		assert.Equal(t, http.StatusBadRequest, actualResponse.Status)
//		assert.Equal(t, constant.EmailAlreadyExists, actualResponse.Message)
//	})
//	t.Run("error while querying DB to check if user exists", func(t *testing.T) {
//		dbErr := errors.New("error while connecting to DB")
//		err := utils.LogError(http.StatusInternalServerError, dbErr.Error(), constant.InternalServerError)
//		mockUserService.Mock.On("UserAlreadyExists", mockUserReq.Email).Return(err, true).Once()
//
//		jsonUserResp, _ := json.Marshal(mockUserReq)
//		mockUserReq, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonUserResp))
//		rec := httptest.NewRecorder()
//		route.ServeHTTP(rec, mockUserReq)
//
//		_ = json.NewDecoder(rec.Body).Decode(&actualResponse)
//
//		assert.Equal(t, http.StatusInternalServerError, rec.Code)
//		assert.Equal(t, http.StatusInternalServerError, actualResponse.Status)
//		assert.Equal(t, constant.InternalServerError, actualResponse.Message)
//	})
//}
