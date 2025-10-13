package internal

import (
	"context"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func noPanicBeHappy(ctx context.Context) {
	if err := recover(); err != nil {
		runtime.EventsEmit(ctx, "runtime:error", err)
	}
}

type CacheManagerApi struct {
	singletonLock sync.Mutex
	seed          *snowflake.Node
	ctx           context.Context
}

func NewCacheManagerApi() *CacheManagerApi {
	var api = &CacheManagerApi{}

	seed, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	api.seed = seed

	return api
}

func (api *CacheManagerApi) Startup(ctx context.Context) {
	api.ctx = ctx
}
