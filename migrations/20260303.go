package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// T20260303 将 user_config 表的 notify JSON 字段从 bool 类型迁移为 int 类型
// true -> 1(开启), false -> 2(关闭), null -> 0(默认值，等同开启)
func T20260303() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20260303",
		Migrate: func(tx *gorm.DB) error {
			fields := []string{
				"create_script", "script_update", "script_issue",
				"script_issue_comment", "score", "at",
				"script_report", "script_report_comment",
			}
			for _, field := range fields {
				if err := tx.Exec(`
					UPDATE cm_user_config SET notify = JSON_SET(
						COALESCE(notify, '{}'),
						CONCAT('$.', ?),
						CASE
							WHEN JSON_EXTRACT(notify, CONCAT('$.', ?)) = true THEN 1
							WHEN JSON_EXTRACT(notify, CONCAT('$.', ?)) = false THEN 2
							ELSE 0
						END
					) WHERE notify IS NOT NULL
				`, field, field, field).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
