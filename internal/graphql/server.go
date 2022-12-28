package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/graphql/graph"
	"github.com/duyike/greddit/internal/graphql/graph/dataloader"
	"github.com/duyike/greddit/internal/graphql/graph/resolver"
	"github.com/duyike/greddit/internal/pkg"
)

var (
	_, _ = zap.NewProduction(zap.Fields(zap.String("type", "graphql")))
)

type App struct {
	pkg.BaseApp
}

func (a *App) Init() (pkg.App, error) {
	_, err := a.BaseApp.Init()
	if err != nil {
		return nil, err
	}
	a.App = initFiberApp()
	return a, nil
}

var _ pkg.App = &App{}

func initFiberApp() *fiber.App {
	fiberApp := fiber.New(fiber.Config{ErrorHandler: middleware.CustomErrorHandler()})
	fiberApp.Use(middleware.UserAuth())

	fiberApp.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("ok")
	})

	loader := dataloader.NewDataLoader()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}}))
	loaderSrv := dataloader.Middleware(loader, srv)

	fiberApp.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))
	fiberApp.Post("/query", adaptor.HTTPHandler(loaderSrv))
	return fiberApp
}
