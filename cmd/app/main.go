package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/api"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/resource_repo"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/task/consumer"
	"github.com/scriptscat/scriptlist/internal/task/crontab"
	"github.com/scriptscat/scriptlist/migrations"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("scriptlist")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	// 注册repository
	script_repo.RegisterScript(script_repo.NewScriptRepo())
	script_repo.RegisterScriptCode(script_repo.NewScriptCodeRepo())

	script_repo.RegisterScriptDomain(script_repo.NewScriptDomainRepo())
	script_repo.RegisterScriptCategory(script_repo.NewScriptCategoryRepo())
	script_repo.RegisterScriptCategoryList(script_repo.NewScriptCategoryListRepo())
	script_repo.RegisterMigrate(script_repo.NewMigrateRepo())
	script_repo.RegisterLibDefinition(script_repo.NewLibDefinitionRepo())
	script_repo.RegisterScriptWatch(script_repo.NewScriptWatchRepo())

	script_repo.RegisterScriptDateStatistics(script_repo.NewScriptDateStatistics())
	script_repo.RegisterScriptStatistics(script_repo.NewScriptStatistics())
	//注册评分
	script_repo.RegisterScriptScore(script_repo.NewScriptScore())

	statistics_repo.RegisterStatistics(statistics_repo.NewStatistics())

	issue_repo.RegisterScriptIssue(issue_repo.NewScriptIssue())
	issue_repo.RegisterScriptIssueComment(issue_repo.NewScriptIssueComment())
	issue_repo.RegisterScriptIssueWatch(issue_repo.NewScriptIssueWatch())

	user_repo.RegisterUser(user_repo.NewUserRepo())
	user_repo.RegisterFollow(user_repo.NewFollowRepo())
	user_repo.RegisterUserConfig(user_repo.NewUserConfig())

	resource_repo.RegisterResource(resource_repo.NewResource())

	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(trace.Trace)).
		Registry(cago.FuncComponent(db.Database)).
		Registry(cago.FuncComponent(redis.Redis)).
		Registry(cago.FuncComponent(cache.Cache)).
		Registry(cago.FuncComponent(elasticsearch.Elasticsearch)).
		Registry(cago.FuncComponent(broker.Broker)).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(db.Default())
		})).
		Registry(cago.FuncComponent(consumer.Consumer)).
		Registry(cago.FuncComponent(crontab.Crontab)).
		RegistryCancel(mux.Http(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
