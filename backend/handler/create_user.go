package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"some-application/backend/contract"
	"some-application/backend/kafka"
	"some-application/backend/kafka/message"
	"some-application/backend/model"
	"some-application/backend/service/user"
	"some-application/backend/utils"
	"some-application/backend/utils/constant"
)

type CreateUserHandler struct {
	userService   user.Service
	kafkaProducer kafka.ClientProducer
}

func NewUserHandler(service user.Service, producer kafka.ClientProducer) CreateUserHandler {
	return CreateUserHandler{
		userService:   service,
		kafkaProducer: producer,
	}
}

func (h *CreateUserHandler) CreateUser(ctx *gin.Context) {
	var req model.User

	validationErr := validateCreateUserRequest(ctx, &req)
	if validationErr != nil {
		log.Println(validationErr.Error)
		ctx.JSON(validationErr.Response.Status, validationErr.Response)
		return
	}

	err, emailExists := h.userService.UserAlreadyExists(req.Email)
	if emailExists {
		log.Println(err.Error)
		ctx.JSON(err.Response.Status, err.Response)
		return
	}

	err = h.userService.HashPassword(&req)
	if err != nil {
		log.Println(err.Error)
		ctx.JSON(err.Response.Status, err.Response)
		return
	}

	err = h.userService.AddUser(&req)
	if err != nil {
		log.Println(err.Error)
		ctx.JSON(err.Response.Status, err.Response)
		return
	}

	ctx.JSON(http.StatusCreated, &utils.Response{
		Message: constant.UserAddedSuccessfully,
		Status:  http.StatusCreated,
	})

	kafkaErr := h.kafkaProducer.EmitEvent(kafka.NewUserEvent, kafka.UserTopic, req.Id.String(), message.NewUser{Email: req.Email})
	if kafkaErr != nil {
		log.Print(kafkaErr)
	}
}

func validateCreateUserRequest(ctx *gin.Context, user *model.User) *utils.Error {
	var req contract.CreateUserReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		var ve validator.ValidationErrors
		var msg string
		if errors.As(err, &ve) {
			msg = utils.ValidationError(ve[0].Tag())
		}
		return utils.LogError(http.StatusBadRequest, err.Error(), ve[0].Field()+": "+msg)
	}
	jsonRes, jsonErr := json.Marshal(req)
	if jsonErr != nil {
		return utils.LogError(http.StatusInternalServerError, jsonErr.Error(), constant.InternalServerError)
	}
	jsonErr = json.Unmarshal(jsonRes, user)
	if jsonErr != nil {
		return utils.LogError(http.StatusInternalServerError, jsonErr.Error(), constant.InternalServerError)
	}
	return nil
}
