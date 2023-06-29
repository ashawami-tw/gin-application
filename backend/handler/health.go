package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (h *HealthHandler) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "healthy")
}
