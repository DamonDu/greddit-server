package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/graphql"
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
		addr     = fmt.Sprintf(":%s", os.Getenv("PORT"))
		shutdown = make(chan struct{})
	)

	app, err := (&graphql.App{}).Init()
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
