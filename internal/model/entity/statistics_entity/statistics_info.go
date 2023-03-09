package statistics_entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Whitelist struct {
	Whitelist []string
}

func (w *Whitelist) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, w)
	return err
}

func (w *Whitelist) Value() (driver.Value, error) {
	ret, err := json.Marshal(w)
	return ret, err
}

type StatisticsInfo struct {
	ID            int64      `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID      int64      `gorm:"column:script_id;type:bigint(20);not null;index:script_id,unique"`
	StatisticsKey string     `gorm:"column:statistics_key;type:varchar(128);index:statistics_key,unique"`
	Whitelist     *Whitelist `gorm:"column:whitelist;type:json;null"`
	Status        int        `gorm:"column:status;type:tinyint(2);not null;default:1"`
	Createtime    int64      `gorm:"column:createtime;type:bigint(20);not null"`
}
