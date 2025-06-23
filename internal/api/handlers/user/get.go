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

func GetUserHandler(c *gin.Context) {
	conn := storage.GetConnector()
	ctx := c.Request.Context()

	telegramIdStr := c.Query("telegram_id")
	discordIdStr := c.Query("discord_id")
	license := c.Query("license")

	params := storage.GetUserParams{}
	var err error

	if telegramIdStr != "" {
		params.TelegramId, err = strconv.Atoi(telegramIdStr)
		if err != nil {
			logger.Debug(ctx, "invalid telegram_id", zap.Error(err))
			c.JSON(api_utils.FormErrResponse(http.StatusBadRequest, "telegram_id must be an integer"))
			return
		}
	} else if discordIdStr != "" {
		params.DiscordId, err = strconv.Atoi(discordIdStr)
		if err != nil {
			logger.Debug(ctx, "invalid discord_id", zap.Error(err))
			c.JSON(api_utils.FormErrResponse(http.StatusBadRequest, "discord_id must be an integer"))
			return
		}
	} else if license != "" {
		params.License = license
	} else {
		c.JSON(api_utils.FormErrResponse(http.StatusBadRequest, "at least telegram_id, discord_id or license must be provided"))
		return
	}

	user, err := conn.GetUser(ctx, params)
	if err != nil {
		logger.Error(ctx, "failed to get user", zap.Error(err))
		c.JSON(api_utils.FormInternalErrResponse())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   user,
	})
}
