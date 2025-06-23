package router

import (
	"github.com/dzhisl/license-api/internal/api/handlers/ping"
	"github.com/dzhisl/license-api/internal/api/handlers/user"
	"github.com/dzhisl/license-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	r.Use(middleware.RequestIDMiddleware())
	RouterGroup := r.Group("/api")

	registerPublicRoutes(RouterGroup)
	registerPrivateRoutes(RouterGroup)
	return r
}

func registerPublicRoutes(r *gin.RouterGroup) {
	r.GET("/ping", ping.PingHandler)
}

func registerPrivateRoutes(r *gin.RouterGroup) {
	r.Use(middleware.AdminAuthMiddleware)
	r.GET("/ping2", ping.PingHandler)
	r.POST("createUser", user.CreateUserHandler)
	r.GET("/getUser", user.GetUserHandler)
}
