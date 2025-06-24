package license

import (
	"net/http"
	"time"

	"github.com/dzhisl/license-api/internal/api/utils"
	"github.com/dzhisl/license-api/internal/storage"
	"github.com/dzhisl/license-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type verifyLicenseRequest struct {
	License string `json:"license" binding:"required"`
	HWID    string `json:"hwid" binding:"required"`
}

// @Summary Verify license
// @Description Verify license by license string and HWID
// @Tags license
// @Accept json
// @Produce json
// @Param request body verifyLicenseRequest true "payload"
// @Router /license/verify [post]
func VerifyLicenseHandler(c *gin.Context) {
	var req verifyLicenseRequest
	ctx := c.Request.Context()
	conn := storage.GetConnector()

	if err := c.ShouldBindJSON(&req); err != nil {
		status, resp := utils.FormInvalidRequestResponse()
		c.JSON(status, resp)
		return
	}

	user, err := conn.GetUser(ctx, storage.GetUserParams{License: req.License})
	if err != nil || user == nil {
		logger.Error(ctx, "failed to find user", zap.Error(err))
		status, resp := utils.FormErrResponse(http.StatusNotFound, "license not found")
		c.JSON(status, resp)
		return
	}

	license := user.License

	if license.Status != storage.Active {
		status, resp := utils.FormErrResponse(http.StatusForbidden, "license not active")
		c.JSON(status, resp)
		return
	}

	// NEW: check if HWID is already registered
	hwidExists := false
	for _, hwid := range license.Devices {
		if hwid == req.HWID {
			hwidExists = true
			break
		}
	}

	if len(license.Devices) >= license.MaxActivations && !hwidExists {
		status, resp := utils.FormErrResponse(http.StatusForbidden, "device limit reached — new device not allowed")
		c.JSON(status, resp)
		return
	}

	if time.Now().Unix() >= int64(license.ExpiresAt) {
		status, resp := utils.FormErrResponse(http.StatusForbidden, "license expired")
		c.JSON(status, resp)
		return
	}

	c.JSON(utils.FormResponse("license is valid"))
}
