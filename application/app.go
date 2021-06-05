package application

import "github.com/damondu/greddit/infrastructure/persistence"

type App struct {
	User *UserApp
	Post *PostApp
}

func NewApp(repositories *persistence.Repositories) *App {
	return &App{
		User: NewUserApp(repositories),
		Post: NewPostApp(repositories),
	}
}
