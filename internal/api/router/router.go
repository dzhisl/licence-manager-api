package router

import (
	"github.com/dzhisl/license-api/internal/api/handlers/license"
	"github.com/dzhisl/license-api/internal/api/handlers/ping"
	"github.com/dzhisl/license-api/internal/api/handlers/user"
	"github.com/dzhisl/license-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	middleware.PrometheusInit()

	r.Use(middleware.RequestIDMiddleware(), middleware.TrackMetrics())
	r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL:                  "/swagger/doc.json", // URL to the generated swagger.json
		DocExpansion:         "none",
		PersistAuthorization: true,
	}, swaggerfiles.Handler))

	RouterGroup := r.Group("/api")

	registerPublicRoutes(*RouterGroup)
	registerPrivateRoutes(*RouterGroup)
	return r
}

func registerPublicRoutes(r gin.RouterGroup) {
	limiter := middleware.NewClientLimiter(1, 5) // 1 req/sec, burst up to 5
	r.Use(middleware.RateLimitMiddleware(limiter))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("ping", ping.PingHandler)
	r.POST("license/verify", license.VerifyLicenseHandler)
}

func registerPrivateRoutes(r gin.RouterGroup) {
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
