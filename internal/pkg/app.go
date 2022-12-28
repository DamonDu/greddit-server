package pkg

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/pkg/db"
	"github.com/duyike/greddit/internal/repository"
	"github.com/duyike/greddit/internal/service"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "app")))
)

type App interface {
	Init() (App, error)
	GracefulShutdown(shutdown chan struct{})
	Listen(addr string) error
	FiberApp() *fiber.App
}

type BaseApp struct {
	*fiber.App
	shutdowns []func() error
}

func (a *BaseApp) Init() (App, error) {
	var (
		shutdowns []func() error
	)
	rand.Seed(time.Now().UnixNano())
	database, err := db.NewDb()
	if err != nil {
		return nil, err
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	shutdowns = append(shutdowns, sqlDB.Close)

	err = autoMigrate(database)
	if err != nil {
		return nil, err
	}

	repository.Init(database)
	service.Init()

	a.App = fiber.New()
	a.shutdowns = shutdowns
	return a, err
}

func (a *BaseApp) GracefulShutdown(shutdown chan struct{}) {
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

func (a *BaseApp) Listen(addr string) error {
	return a.App.Listen(addr)
}

func (a *BaseApp) FiberApp() *fiber.App {
	return a.App
}

var _ App = &BaseApp{}

func autoMigrate(database *gorm.DB) error {
	models := []interface{}{&model.User{}, &model.Post{}}
	for _, m := range models {
		err := database.AutoMigrate(m)
		if err != nil {
			return err
		}
	}
	return nil
}
