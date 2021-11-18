package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/api"
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
		port     = os.Getenv("PORT")
		addr     = ":" + port
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
