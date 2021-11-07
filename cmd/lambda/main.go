package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	adapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/api"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "lambda")))
)

func main() {
	app, err := api.NewApp()
	if err != nil {
		logger.Fatal("new app error", zap.Error(err))
	}

	lambdaApp := adapter.New(app.FiberApp())
	lambda.Start(lambdaApp.ProxyWithContext)
}
