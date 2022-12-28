package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/duyike/greddit/internal/api/handler"
	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/pkg"
)

var (
	_, _ = zap.NewProduction(zap.Fields(zap.String("type", "api")))
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
	fiberApp.Use(fiberRecover.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowHeaders: strings.Join([]string{
			fiber.HeaderOrigin,
			fiber.HeaderContentLength,
			fiber.HeaderContentType,
			fiber.HeaderAuthorization,
		}, ","),
		AllowCredentials: true,
	}))
	fiberApp.Use(fiberLogger.New())
	fiberApp.Use(middleware.UserAuth())

	fiberApp.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("ok")
	})
	fiberApp.Mount("/user", handler.NewUserHandler().App)
	fiberApp.Mount("/post", handler.NewPostHandler().App)
	return fiberApp
}
