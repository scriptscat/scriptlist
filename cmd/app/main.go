package main

import (
	"context"
	"log"

	"github.com/scriptscat/scriptlist/internal/repository/feedback_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_profile_repo"

	"github.com/cago-frame/cago"
	"github.com/cago-frame/cago/configs"
	"github.com/cago-frame/cago/database/cache"
	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/database/elasticsearch"
	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/broker"
	"github.com/cago-frame/cago/pkg/component"
	"github.com/cago-frame/cago/server/cron"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/api"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/notification_repo"
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
	script_repo.RegisterScriptAccess(script_repo.NewScriptAccess())
	script_repo.RegisterScriptGroup(script_repo.NewScriptGroup())
	script_repo.RegisterScriptGroupMember(script_repo.NewScriptGroupMember())
	script_repo.RegisterScriptInvite(script_repo.NewScriptInvite())
	//注册评分
	script_repo.RegisterScriptScore(script_repo.NewScriptScore())
	// 收藏夹
	script_repo.RegisterScriptFavorite(script_repo.NewScriptFavorite())
	script_repo.RegisterScriptFavoriteFolder(script_repo.NewScriptFavoriteFolder())

	statistics_repo.RegisterScriptStatistics(statistics_repo.NewScriptStatistics())
	statistics_repo.RegisterStatisticsInfo(statistics_repo.NewStatisticsInfo())

	issue_repo.RegisterScriptIssue(issue_repo.NewScriptIssue())
	issue_repo.RegisterScriptIssueComment(issue_repo.NewScriptIssueComment())
	issue_repo.RegisterScriptIssueWatch(issue_repo.NewScriptIssueWatch())

	user_repo.RegisterUser(user_repo.NewUserRepo())
	user_repo.RegisterFollow(user_repo.NewFollowRepo())
	user_repo.RegisterUserConfig(user_repo.NewUserConfig())
	user_profile_repo.RegisterUserProfile(user_profile_repo.NewUserProfile())

	feedback_repo.RegisterFeedback(feedback_repo.NewFeedback())

	resource_repo.RegisterResource(resource_repo.NewResource())

	notification_repo.RegisterNotification(notification_repo.NewNotificationRepo())

	err = cago.New(ctx, cfg).
		Registry(component.Core()).
		Registry(db.Database()).
		Registry(cago.FuncComponent(redis.Redis)).
		Registry(cache.Cache()).
		Registry(cago.FuncComponent(elasticsearch.Elasticsearch)).
		Registry(cago.FuncComponent(broker.Broker)).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(db.Default())
		})).
		Registry(cago.FuncComponent(consumer.Consumer)).
		Registry(cron.Cron()).
		Registry(cago.FuncComponent(crontab.Crontab)).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
