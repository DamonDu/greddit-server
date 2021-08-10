package main

import (
	"github.com/damondu/greddit/internal/api"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
)

func main() {
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
