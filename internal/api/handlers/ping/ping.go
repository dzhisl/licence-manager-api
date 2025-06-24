package ping

import (
	"github.com/dzhisl/license-api/internal/api/utils"
	"github.com/dzhisl/license-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// @Summary Simple ping endpoint
// @Description Ping pong
// @Tags user
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	logger.Info(c.Request.Context(), "ponged")
	c.JSON(utils.FormResponse("pong"))
}
