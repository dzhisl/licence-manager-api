package user

import (
	"net/http"
	"strconv"

	api_utils "github.com/dzhisl/license-api/internal/api/utils"
	"github.com/dzhisl/license-api/internal/storage"
	"github.com/dzhisl/license-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type addDeviceRequest struct {
	HWID string `json:"hwid" binding:"required"`
}

// @Summary Bind HWID to user
// @Description Add a device HWID to user device list
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body addDeviceRequest true "payload"
// @Success 200 {object} statusResponse
// @Failure 400 {object} invalidBodyErrResponse
// @Failure 500 {object} internalErrResponse
// @Security ApiKeyAuth
// @Router /user/{user_id}/device [post]
func AddDeviceHandler(c *gin.Context) {
	conn := storage.GetConnector()
	ctx := c.Request.Context()

	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logger.Debug(ctx, "invalid user_id", zap.Error(err))
		c.JSON(api_utils.FormErrResponse(http.StatusBadRequest, "user_id must be an integer"))
		return
	}

	var req addDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(api_utils.FormInvalidRequestResponse())
		return
	}

	err = conn.AddHwidSession(ctx, userId, req.HWID)
	if err != nil {
		logger.Error(ctx, "failed to add hwid session", zap.Error(err))
		c.JSON(api_utils.FormInternalErrResponse())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
