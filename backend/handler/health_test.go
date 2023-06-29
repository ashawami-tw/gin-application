package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	healthHandler := NewHealthHandler()

	req, _ := http.NewRequest("GET", "/health", nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	healthHandler.HealthHandler(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)
}
