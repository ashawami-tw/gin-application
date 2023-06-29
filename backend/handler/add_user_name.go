package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"some-application/backend/contract"
	"some-application/backend/kafka"
	"some-application/backend/model"
	"some-application/backend/service/user"
	"some-application/backend/service/username"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
)

type AddUserNameHandler struct {
	userService     user.Service
	usernameService username.Service
	kafkaProducer   kafka.ClientProducer
}

func NewAddUserNameHandler(service user.Service, usernameService username.Service, producer kafka.ClientProducer) AddUserNameHandler {
	return AddUserNameHandler{
		userService:     service,
		usernameService: usernameService,
		kafkaProducer:   producer,
	}
}

func (h *AddUserNameHandler) AddUserName(ctx *gin.Context) {
	var userName model.UserName
	validationErr := validateAddUserNameRequest(ctx, &userName)

	if validationErr != nil {
		log.Println(validationErr.Error)
		ctx.JSON(validationErr.Response.Status, validationErr.Response)
		return
	}

	err, userIdExists := h.userService.UserIdExists(userName.UserId)
	if !userIdExists {
		log.Println(err.Error)
		ctx.JSON(err.Response.Status, err.Response)
		return
	}

	err = h.usernameService.AddUserName(&userName)
	if err != nil {
		log.Println(err.Error)
		ctx.JSON(err.Response.Status, err.Response)
		return
	}

	ctx.JSON(http.StatusOK, userName)
}

func validateAddUserNameRequest(ctx *gin.Context, username *model.UserName) *utils.Error {
	var req contract.AddUserNameReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return utils.LogError(http.StatusBadRequest, err.Error(), constant.ValidationError)
	}
	jsonRes, jsonErr := json.Marshal(&req)
	if jsonErr != nil {
		return utils.LogError(http.StatusInternalServerError, jsonErr.Error(), constant.InternalServerError)
	}
	jsonErr = json.Unmarshal(jsonRes, username)
	if jsonErr != nil {
		return utils.LogError(http.StatusInternalServerError, jsonErr.Error(), constant.InternalServerError)
	}
	return nil
}
