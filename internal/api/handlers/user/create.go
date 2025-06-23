package user

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	api_utils "github.com/dzhisl/license-api/internal/api/utils"
	"github.com/dzhisl/license-api/internal/storage"
	"github.com/dzhisl/license-api/pkg/logger"
	"github.com/dzhisl/license-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createUserRequest struct {
	TelegramId     int `json:"telegram_id"`
	DiscordId      int `json:"discord_id"`
	MaxActivations int `json:"max_activations" binding:"required"`
	Expiration     int `json:"expires_at"  binding:"required"`
}

func CreateUserHandler(c *gin.Context) {
	conn := storage.GetConnector()
	ctx := c.Request.Context()

	var reqBody createUserRequest
	if err := c.ShouldBind(&reqBody); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(api_utils.FormInvalidRequestResponse())
		return
	}

	now := time.Now().Unix()
	user := storage.User{
		Id:        genUserID(),
		CreatedAt: storage.Timestamp(now),
		License: storage.License{
			Key:            utils.GenLicense(),
			MaxActivations: reqBody.MaxActivations,
			IssuedAt:       storage.Timestamp(now),
			ExpiresAt:      storage.Timestamp(reqBody.Expiration),
			Status:         storage.Active,
		},
	}
	switch {
	case reqBody.TelegramId != 0:
		user.TelegramId = reqBody.TelegramId
	case reqBody.DiscordId != 0:
		user.DiscordId = reqBody.DiscordId
	default:
		c.JSON(api_utils.FormErrResponse(http.StatusBadRequest, "at least discord or telegram ID must be provided"))
		return
	}

	if err := conn.CreateUser(ctx, user); err != nil {
		logger.Error(ctx, err.Error(), zap.Any("request_body", reqBody))
		c.JSON(api_utils.FormInternalErrResponse())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   user,
	})
}

func genUserID() int {
	const digits = "123456789"

	var sb strings.Builder
	for i := 0; i < 8; i++ {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}
	userId, _ := strconv.Atoi(sb.String())
	return userId
}
