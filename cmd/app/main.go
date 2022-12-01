package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/api"
	"github.com/scriptscat/scriptlist/internal/pkg/consumer"
	"github.com/scriptscat/scriptlist/migrations"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("scriptlist")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(trace.Trace)).
		Registry(cago.FuncComponent(db.Database)).
		Registry(cago.FuncComponent(redis.Redis)).
		Registry(cago.FuncComponent(cache.Cache)).
		Registry(cago.FuncComponent(broker.Broker)).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(db.Default())
		})).
		Registry(cago.FuncComponent(consumer.Consumer)).
		RegistryCancel(mux.Http(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
