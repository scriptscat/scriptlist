package script_repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/pkg/utils/httputils"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptRepo interface {
	Find(ctx context.Context, id int64) (*entity.Script, error)
	Create(ctx context.Context, script *entity.Script) error
	Update(ctx context.Context, script *entity.Script) error
	Delete(ctx context.Context, id int64) error

	Search(ctx context.Context, keyword, sort string, scriptType int, page httputils.PageRequest) ([]*entity.ScriptSearch, int64, error)
}

var defaultScript ScriptRepo

func Script() ScriptRepo {
	return defaultScript
}

func RegisterScript(i ScriptRepo) {
	defaultScript = i
}

type scriptRepo struct {
}

func NewScriptRepo() ScriptRepo {
	return &scriptRepo{}
}

func (u *scriptRepo) Find(ctx context.Context, id int64) (*entity.Script, error) {
	ret := &entity.Script{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptRepo) Create(ctx context.Context, script *entity.Script) error {
	return db.Ctx(ctx).Create(script).Error
}

func (u *scriptRepo) Update(ctx context.Context, script *entity.Script) error {
	return db.Ctx(ctx).Updates(script).Error
}

func (u *scriptRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Script{ID: id}).Error
}

func (u *scriptRepo) Search(ctx context.Context, keyword, sort string, scriptType int, page httputils.PageRequest) ([]*entity.ScriptSearch, int64, error) {
	if keyword != "" {
		// 暂时不支持排序等
		return u.SearchByEs(ctx, keyword, page)
	}
	// 无关键字从mysql数据库中获取
	return nil, 0, nil
}

// SearchByEs 通过elasticsearch搜索
func (u *scriptRepo) SearchByEs(ctx context.Context, keyword string, page httputils.PageRequest) ([]*entity.ScriptSearch, int64, error) {
	script := &entity.ScriptSearch{}
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
	m := &elasticsearch.SearchResponse[*entity.ScriptSearch]{}
	if err := json.Unmarshal(respByte, &m); err != nil {
		return nil, 0, err
	}
	ret := make([]*entity.ScriptSearch, 0)
	for _, v := range m.Hits.Hits {
		ret = append(ret, v.Source)
	}
	return ret, m.Hits.Total.Value, nil
}
