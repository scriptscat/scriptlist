package script

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	script2 "github.com/scriptscat/scriptlist/internal/repository/script_repo"
)

type script struct {
}

func NewScript() script2.IScript {
	return &script{}
}

func (u *script) Find(ctx context.Context, id int64) (*entity.Script, error) {
	ret := &entity.Script{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *script) Create(ctx context.Context, script *entity.Script) error {
	return db.Ctx(ctx).Create(script).Error
}

func (u *script) Update(ctx context.Context, script *entity.Script) error {
	return db.Ctx(ctx).Updates(script).Error
}

func (u *script) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Script{ID: id}).Error
}

func (u *script) Search(ctx context.Context, keyword, sort string, scriptType int, page httputils.PageRequest) ([]*model.ScriptSearch, int64, error) {
	if keyword != "" {
		// 暂时不支持排序等
		return u.SearchByEs(ctx, keyword, page)
	}
	// 无关键字从mysql数据库中获取
	return nil, 0, nil
}

// SearchByEs 通过elasticsearch搜索
func (u *script) SearchByEs(ctx context.Context, keyword string, page httputils.PageRequest) ([]*model.ScriptSearch, int64, error) {
	script := &model.ScriptSearch{}
	search := elasticsearch.Ctx(ctx).Search
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": keyword,
			},
		},
		"size": page.Limit,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}
	resp, err := elasticsearch.Ctx(ctx).Search(
		search.WithIndex(script.CollectionName()),
		search.WithBody(&buf),
		search.WithTrackTotalHits(true),
		search.WithPretty())
	if err != nil {
		return nil, 0, err
	}
	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	if resp.IsError() {
		return nil, 0, fmt.Errorf("elasticsearch error: [%s] %s", resp.Status(), respByte)
	}
	m := &elasticsearch.SearchResponse[*model.ScriptSearch]{}
	if err := json.Unmarshal(respByte, &m); err != nil {
		return nil, 0, err
	}
	ret := make([]*model.ScriptSearch, 0)
	for _, v := range m.Hits.Hits {
		ret = append(ret, v.Source)
	}
	return ret, m.Hits.Total.Value, nil
}
