package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/golang/glog"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	entity2 "github.com/scriptscat/scriptweb/internal/domain/statistics/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"gorm.io/gorm"
)

func Migrate() error {
	m := gormigrate.New(db.Db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1617420365",
			Migrate: func(tx *gorm.DB) error {
				if !tx.Migrator().HasTable(&entity.Script{}) {
					if err := tx.AutoMigrate(&entity.Script{}); err != nil {
						return err
					}
				}
				if !tx.Migrator().HasTable(&entity.ScriptCode{}) {
					if err := tx.AutoMigrate(&entity.ScriptCode{}); err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&entity.Script{})
			},
		}, {
			ID: "1622952090",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&entity.ScriptCategoryList{}); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity.ScriptCategory{}); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity.ScriptScore{}); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity2.StatisticsDownload{}); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity2.StatisticsUpdate{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable(&entity2.StatisticsDownload{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&entity2.StatisticsUpdate{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&entity.ScriptScore{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&entity.ScriptCategoryList{}); err != nil {
					return err
				}
				return tx.Migrator().DropTable(&entity.ScriptCategory{})
			},
		}, {
			ID: "1627371597",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&entity.ScriptCode{}, "MetaJson"); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity.ScriptDomain{}); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity.ScriptDateStatistics{}); err != nil {
					return err
				}
				if err := tx.AutoMigrate(&entity.ScriptStatistics{}); err != nil {
					return err
				}
				// 处理meta_json和domain信息
				go func() {
					if err := DealMetaInfo(); err != nil {
						glog.Fatal("deal meta info: %v", err)
					}
				}()
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&entity.ScriptCode{}, "MetaJson"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&entity.ScriptDomain{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&entity.ScriptDateStatistics{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&entity.ScriptStatistics{}); err != nil {
					return err
				}
				return nil
			},
		}, {
			ID: "1627723150",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&entity.ScriptStatistics{}, "score"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.ScriptStatistics{}, "score_count"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&entity.ScriptStatistics{}, "score"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.ScriptStatistics{}, "score_count"); err != nil {
					return err
				}
				return nil
			},
		}, {
			ID: "1627908382",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&entity.Script{}, "type"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.Script{}, "public"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.Script{}, "unwell"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.Script{}, "sync_url"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.Script{}, "sync_mode"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.Script{}, "content_url"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&entity.Script{}, "definition_url"); err != nil {
					return err
				}
				return tx.Migrator().AutoMigrate(&entity.LibDefinition{})
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&entity.Script{}, "type"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.Script{}, "public"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.Script{}, "unwell"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.Script{}, "sync_url"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.Script{}, "sync_mode"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.Script{}, "content_url"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&entity.Script{}, "definition_url"); err != nil {
					return err
				}
				return tx.Migrator().DropTable(&entity.LibDefinition{})
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
