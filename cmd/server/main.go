package main

import (
	"context"

	_ "github.com/dzhisl/license-api/docs"
	"github.com/dzhisl/license-api/internal/api/router"
	"github.com/dzhisl/license-api/internal/storage"
	"github.com/dzhisl/license-api/pkg/config"
	"github.com/dzhisl/license-api/pkg/logger"
)

func initApp(ctx context.Context) {
	config.InitConfig()
	logger.InitLogger()
	storage.InitStorage(ctx)
}

// @title License Manager API
// @version 1.0
// @description API for managing user licenses.
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	ctx := context.TODO()
	initApp(ctx)
	r := router.InitRouter()
	logger.Info(ctx, "running API")
	r.Run(":8080")
}
