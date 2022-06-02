package persistence

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	goRedis "github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	persistence6 "github.com/scriptscat/scriptlist/internal/service/issue/infrastructure/persistence"
	persistence2 "github.com/scriptscat/scriptlist/internal/service/resource/infrastructure/persistence"
	persistence5 "github.com/scriptscat/scriptlist/internal/service/safe/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/internal/service/script/infrastructure/persistence"
	persistence3 "github.com/scriptscat/scriptlist/internal/service/statistics/infrastructure/persistence"
	persistence4 "github.com/scriptscat/scriptlist/internal/service/user/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/migrations"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repositories struct {
	Db         *gorm.DB
	Redis      *goRedis.Client
	Cache      cache.Cache
	Script     *persistence.Repositories
	Statistics *persistence3.StatisRepositories
	Resource   *persistence2.Repositories
	User       *persistence4.UserRepositories
	Safe       *persistence5.SafeRepositories
	Issue      *persistence6.IssueRepositories
}

func NewRepositories(db *gorm.DB, redis *goRedis.Client, cache cache.Cache) *Repositories {
	return &Repositories{
		Db:         db,
		Redis:      redis,
		Cache:      cache,
		Script:     persistence.NewRepositories(db, redis, cache),
		Resource:   persistence2.NewRepositories(redis),
		Statistics: persistence3.NewRepositories(db, redis),
		User:       persistence4.NewRepositories(db, redis, cache),
		Safe:       persistence5.NewRepositories(redis),
		Issue:      persistence6.NewRepositories(db, redis),
	}
}

func (r *Repositories) Migrations() error {
	return migrations.RunMigrations(r.Db)
}

func (r *Repositories) MockDB() (*Repositories, *sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	gdb, _ := gorm.Open(mysql.Open(config.AppConfig.Mysql.Dsn), &gorm.Config{})
	return &Repositories{
		Db: gdb,
	}, db, mock
}
