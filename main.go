package main

import (
	"github.com/damondu/greddit/application"
	"github.com/damondu/greddit/infrastructure/persistence"
	"github.com/damondu/greddit/interface/handler"
	"github.com/damondu/greddit/interface/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	repositories, repoErr := persistence.NewRepositories()
	if repoErr != nil {
		log.Fatal(repoErr)
	}
	apps := application.NewApp(repositories)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))
	router.Use(middleware.AuthToken())
	router.Use(middleware.Errors())

	// user
	userHandler := handler.NewUserHandler(apps)
	userRouter := router.Group("/user")
	{
		userRouter.POST("/me", middleware.AuthUser(), userHandler.Me)
		userRouter.POST("/register", userHandler.Register)
		userRouter.POST("/login", userHandler.Login)
	}

	// post
	postHandler := handler.NewPostHandler(apps)
	postRouter := router.Group("/post")
	{
		postRouter.POST("/pageQuery", postHandler.PageQuery)
	}
	runErr := router.Run(":8080")
	if runErr != nil {
		log.Fatal(runErr)
	}
}
