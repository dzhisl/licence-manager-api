package router

import (
	"github.com/dzhisl/license-api/internal/api/handlers/license"
	"github.com/dzhisl/license-api/internal/api/handlers/ping"
	"github.com/dzhisl/license-api/internal/api/handlers/user"
	"github.com/dzhisl/license-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	r.Use(middleware.RequestIDMiddleware())
	// Register Swagger UI endpoint
	// Note: `PersistAuthorization: true` ensures that the X-API-Key header is preserved after page refreshes.
	r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL:                  "/swagger/doc.json", // URL to the generated swagger.json
		DocExpansion:         "none",
		PersistAuthorization: true,
	}, swaggerfiles.Handler))

	RouterGroup := r.Group("/api")

	registerPublicRoutes(RouterGroup)
	registerPrivateRoutes(RouterGroup)
	return r
}

func registerPublicRoutes(r *gin.RouterGroup) {

	r.GET("ping", ping.PingHandler)
	r.POST("license/verify", license.VerifyLicenseHandler)
}

func registerPrivateRoutes(r *gin.RouterGroup) {
	r.Use(middleware.AdminAuthMiddleware)
	r.POST("user/create", user.CreateUserHandler)
	r.GET("user", user.GetUserHandler)
	r.POST("user/:user_id/device", user.AddDeviceHandler)
	r.DELETE("user/:user_id/device", user.RemoveDeviceHandler)
	r.POST("user/:user_id/devices/reset", user.ResetDevicesHandler)
	r.POST("user/:user_id/license/status", user.ChangeLicenseStatusHandler)
	r.POST("user/:user_id/license/hwid_limit", user.UpdateHwidLimitHandler)
	r.POST("user/:user_id/license/renew", user.RenewLicenseHandler)
	r.POST("user/:user_id/discord", user.BindDiscordHandler)
	r.POST("user/:user_id/telegram", user.BindTelegramHandler)
	r.DELETE("user/:user_id", user.DeleteUserHandler)
}
