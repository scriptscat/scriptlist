package migrations

import (
	"github.com/codfrm/cago/database/clickhouse"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
	"gorm.io/gorm"
)

func T20230228() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230228",
		Migrate: func(tx *gorm.DB) error {
			if err := clickhouse.Default().
				Set("gorm:table_options",
					"ENGINE=ReplacingMergeTree() PARTITION BY script_id "+
						"ORDER BY visitor_id",
				).
				AutoMigrate(&statistics_entity.StatisticsVisitor{}); err != nil {
				return err
			}
			return clickhouse.Default().
				Set("gorm:table_options",
					"ENGINE=ReplacingMergeTree() PARTITION BY script_id "+
						"ORDER BY session_id",
				).AutoMigrate(&statistics_entity.StatisticsCollect{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := clickhouse.Default().Migrator().DropTable(&statistics_entity.StatisticsVisitor{}); err != nil {
				return err
			}
			return clickhouse.Default().Migrator().DropTable(&statistics_entity.StatisticsCollect{})
		},
	}
}
