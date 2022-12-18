package dataloader

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graph-gophers/dataloader"

	"github.com/duyike/greddit/internal/graphql/graph/model"
	dbModel "github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/service"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// DataLoader offers data loaders scoped to a context
type DataLoader struct {
	userLoader *dataloader.Loader
}

// GetUser wraps the User dataloader for efficient retrieval by user ID
func (i *DataLoader) GetUser(ctx context.Context, uid int64) (*model.User, error) {
	thunk := i.userLoader.Load(ctx, Int64Key(uid))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader() *DataLoader {
	return &DataLoader{
		userLoader: dataloader.NewBatchedLoader((&userLoader{}).get),
	}
}

// Middleware injects a DataLoader into the request context so it can be
// used later in the schema resolvers
func Middleware(loader *DataLoader, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loader)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *DataLoader {
	return ctx.Value(loadersKey).(*DataLoader)
}

// userLoader wraps storage and provides a "get" method for the user dataloader
type userLoader struct {
}

// get implements the dataloader for finding many users by Id and returns
// them in the order requested
func (u *userLoader) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	uidList := make([]int64, len(keys))
	for i, key := range keys {
		uidList[i] = int64(key.Raw().(Int64Key))
	}
	// search for those users
	users, err := service.User.BatchGetByUid(uidList)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}
	userMap := users.GroupByInt64((*dbModel.User).GetUserUid)
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for i, uid := range uidList {
		dbUser, ok := userMap[uid]
		if ok {
			results[i] = &dataloader.Result{Data: &model.User{
				UID:      dbUser.Uid,
				Username: dbUser.Username,
				Email:    dbUser.Email,
			}, Error: nil}
		} else {
			results[i] = &dataloader.Result{Data: nil, Error: fmt.Errorf("dbUser not found %d", uid)}
		}
	}
	return results
}
