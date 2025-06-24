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

type renewLicenseRequest struct {
	ExpiresAt int64 `json:"expires_at" binding:"required"`
}

// @Summary Renew license
// @Description Renew the license for a user by user_id
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body renewLicenseRequest true "payload"
// @Success 200 {object} statusResponse
// @Failure 400 {object} invalidBodyErrResponse
// @Failure 500 {object} internalErrResponse
// @Security ApiKeyAuth
// @Router /user/{user_id}/license/renew [post]
func RenewLicenseHandler(c *gin.Context) {
	conn := storage.GetConnector()
	ctx := c.Request.Context()

	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logger.Debug(ctx, "invalid user_id", zap.Error(err))
		c.JSON(api_utils.FormErrResponse(http.StatusBadRequest, "user_id must be an integer"))
		return
	}

	var req renewLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(api_utils.FormInvalidRequestResponse())
		return
	}

	err = conn.RenewLicense(ctx, userId, storage.Timestamp(req.ExpiresAt))
	if err != nil {
		logger.Error(ctx, "failed to renew license", zap.Error(err))
		c.JSON(api_utils.FormInternalErrResponse())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
