package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"some-application/backend/contract"
	"some-application/backend/kafka"
	"some-application/backend/service/user"
	"some-application/backend/service/username"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
	"testing"
)

//func TestNewAddUserNameHandler_Success(t *testing.T) {
//	userId := uuid.New()
//	type args struct {
//		err                error
//		expectedStatusCode int
//		message            string
//		userIdExists       bool
//		mockUsernameReq    *contract.AddUserNameReq
//		mockUsernameResp   *model.UserName
//	}
//
//	successfulUserCreation := args{
//		err:                nil,
//		expectedStatusCode: http.StatusOK,
//		message:            "",
//		userIdExists:       false,
//		mockUsernameReq: &contract.AddUserNameReq{
//			UserId:    userId.String(),
//			FirstName: "Steve",
//			LastName:  "Smith",
//		},
//		mockUsernameResp: &model.UserName{
//			UserId:    userId,
//			FirstName: "Steve",
//			LastName:  "Smith",
//		},
//	}
//}

func TestNewAddUserNameHandler_Validation(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		message            string
		mockUsernameReq    *contract.AddUserNameReq
	}{
		{
			name:               "Empty request body",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq:    &contract.AddUserNameReq{},
		},
		{
			name:               "Empty user id",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    "",
				FirstName: "Steve",
				LastName:  "Smith",
			},
		},
		{
			name:               "Empty firstname",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    uuid.New().String(),
				FirstName: "",
				LastName:  "Smith",
			},
		},
		{
			name:               "Empty lastname",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    uuid.New().String(),
				FirstName: "Steve",
				LastName:  "",
			},
		},
		{
			name:               "Firstname of length 1",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    uuid.New().String(),
				FirstName: "S",
				LastName:  "Smith",
			},
		},
		{
			name:               "Lastname of length 1",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    uuid.New().String(),
				FirstName: "Steve",
				LastName:  "S",
			},
		},
		{
			name:               "Firstname of length 51",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    uuid.New().String(),
				FirstName: "SteveSteveSteveSteveSteveSteveSteveSteveSteveSteveS",
				LastName:  "Smith",
			},
		},
		{
			name:               "Lastname of length 51",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    uuid.New().String(),
				FirstName: "Steve",
				LastName:  "SmithSmithSmithSmithSmithSmithSmithSmithSmithSmithS",
			},
		},
		{
			name:               "User id is not having UUID format",
			expectedStatusCode: http.StatusBadRequest,
			message:            constant.ValidationError,
			mockUsernameReq: &contract.AddUserNameReq{
				UserId:    "UUID",
				FirstName: "Steve",
				LastName:  "Smith",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// dependency
			mockUserService := user.NewUserServiceMock()
			mockUsernameService := username.NewUsernameServiceMock()
			mockProducer := kafka.NewMockProducer()
			usernameHandler := NewAddUserNameHandler(mockUserService, mockUsernameService, mockProducer)

			// prepare user req
			jsonBody, _ := json.Marshal(test.mockUsernameReq)
			req, _ := http.NewRequest("POST", "/username", bytes.NewBuffer(jsonBody))

			// setup to call handler
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = req
			usernameHandler.AddUserName(ctx)

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
