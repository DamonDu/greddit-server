package api

import (
	"github.com/damondu/greddit/internal/api/handler"
	db2 "github.com/damondu/greddit/internal/pkg/db"
	"github.com/damondu/greddit/internal/post"
	"github.com/damondu/greddit/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	logger2 "github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "api")))
)

type App interface {
	GracefulShutdown(shutdown chan struct{})
	Listen(addr string) error
}

type app struct {
	*fiber.App
	shutdowns []func() error
}

func NewApp() (App, error) {
	var (
		shutdowns []func() error
	)
	rand.Seed(time.Now().UnixNano())
	db, err := db2.NewDb()
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	shutdowns = append(shutdowns, sqlDB.Close)

	userRepository := user.NewRepository(db)
	postRepository := post.NewRepository(db)

	userApp := user.NewApp(userRepository)
	postApp := post.NewApp(postRepository, userApp)

	fiberApp := fiber.New()
	fiberApp.Use(recover2.New())
	fiberApp.Use(cors.New(cors.Config{

		AllowHeaders: strings.Join([]string{
			fiber.HeaderOrigin,
			fiber.HeaderContentLength,
			fiber.HeaderContentType,
		}, ","),
		AllowCredentials: true,
	}))
	fiberApp.Use(logger2.New())
	fiberApp.Mount("/user", handler.NewUserHandler(userApp).App)
	fiberApp.Mount("/post", handler.NewPostHandler(postApp).App)
	return app{
		App:       fiberApp,
		shutdowns: shutdowns,
	}, nil
}

func (a app) GracefulShutdown(shutdown chan struct{}) {
	var (
		sigint = make(chan os.Signal, 1)
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	logger.Info("shutting down server gracefully")
	if err := a.Shutdown(); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}
	for i := range a.shutdowns {
		err := a.shutdowns[i]()
		if err != nil {
			logger.Error("sub shutdown error", zap.Any("shutdowns", a.shutdowns[i]))
		}
	}
	close(shutdown)
}

var _ App = &app{}
