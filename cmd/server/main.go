package main

import (
	"context"

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

func main() {
	ctx := context.TODO()
	initApp(ctx)
	r := router.InitRouter()
	logger.Info(ctx, "running API")
	r.Run(":8080")
}
