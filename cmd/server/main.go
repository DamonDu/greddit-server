package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/api"
	"github.com/duyike/greddit/internal/pkg/constant"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "server")))
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		logger.Fatal("Error loading env file")
	}

	var (
		addr     = fmt.Sprintf(":%d", constant.Port)
		shutdown = make(chan struct{})
	)

	app, err := api.NewApp()
	if err != nil {
		logger.Fatal("new app error", zap.Error(err))
	}

	go app.GracefulShutdown(shutdown)
	logger.Info("server start at " + addr)
	err = app.Listen(addr)
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}
	<-shutdown
}
