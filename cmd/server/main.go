package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/api"
	"github.com/duyike/greddit/internal/graphql"
	"github.com/duyike/greddit/internal/pkg"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "server")))
)

func main() {
	cliApp := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Value:   "rest",
				Usage:   "type of API server: rest (default) or graphql",
			},
		},
		Action: runApp,
	}
	_ = cliApp.Run(os.Args)
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runApp(cCtx *cli.Context) error {
	err := godotenv.Load(".env.local")
	if err != nil {
		logger.Fatal("Error loading env file")
		return err
	}

	var (
		addr     = fmt.Sprintf(":%s", os.Getenv("PORT"))
		shutdown = make(chan struct{})
	)

	var app pkg.App
	if cCtx.String("type") == "graphql" {
		app, err = (&graphql.App{}).Init()
	} else {
		app, err = (&api.App{}).Init()
	}
	if err != nil {
		logger.Fatal("new app error", zap.Error(err))
		return err
	}

	go app.GracefulShutdown(shutdown)
	logger.Info("server start at " + addr)
	err = app.Listen(addr)
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
		return err
	}
	<-shutdown
	return nil
}
