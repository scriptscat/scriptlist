package script_entity

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/codfrm/cago/configs"
)

/*
// 索引模板
PUT _index_template/scriptlist.script
{
  "template": {
    "mappings": {
      "properties": {
        "content": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_smart"
        },
        "description": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_smart"
        },
        "name": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_smart"
        }
      }
    }
  },
  "index_patterns": [
    "dev.script"
  ]
}
*/

// ScriptSearch es中储存的数据
type ScriptSearch struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Content       string        `json:"content"`
	Version       string        `json:"version"`
	Changelog     string        `json:"changelog"`
	TotalDownload int64         `json:"total_download"`
	TodayDownload int64         `json:"today_download"`
	Score         float64       `json:"score"`
	Category      []int64       `json:"category"`
	Domain        []string      `json:"domain"`
	Public        Public        `json:"public"`
	Unwell        UnwellContent `json:"unwell"`
	Archive       ScriptArchive `json:"archive"`
	Type          Type          `json:"type"`
	Status        int64         `json:"status"`
	Createtime    int64         `json:"createtime"`
	Updatetime    int64         `json:"updatetime"`
}

func (s *ScriptSearch) CollectionName() string {
	return string(configs.Default().Env) + ".script"
}

func (s *ScriptSearch) Reader() (io.Reader, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
