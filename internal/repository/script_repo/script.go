package script_repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/codfrm/cago/database/cache"
	cache2 "github.com/codfrm/cago/database/cache/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SearchOptions struct {
	Keyword  string
	Type     int
	Sort     string
	UserID   int64
	Self     bool
	Category []int64
	Domain   string
}

//go:generate mockgen -source=script.go -destination=mock/script.go
type ScriptRepo interface {
	Find(ctx context.Context, id int64) (*entity.Script, error)
	Create(ctx context.Context, script *entity.Script) error
	Update(ctx context.Context, script *entity.Script) error
	Delete(ctx context.Context, id int64) error

	Search(ctx context.Context, options *SearchOptions, page httputils.PageRequest) ([]*entity.Script, int64, error)
	// FindSyncScript 查找需要自动同步的脚本
	FindSyncScript(ctx context.Context, page httputils.PageRequest) ([]*entity.Script, error)
	FindSyncPrefix(ctx context.Context, uid int64, prefix string) ([]*entity.Script, error)
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

func (u *scriptRepo) key(id int64) string {
	return "script:" + strconv.FormatInt(id, 10)
}

func (u *scriptRepo) KeyDepend(id int64) *cache2.KeyDepend {
	return cache2.NewKeyDepend(cache.Default(), u.key(id)+":dep")
}

func (u *scriptRepo) Find(ctx context.Context, id int64) (*entity.Script, error) {
	ret := &entity.Script{}
	if err := db.Ctx(ctx).Where("id=?", id).First(ret).Error; err != nil {
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
	return db.Ctx(ctx).Model(&entity.Script{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *scriptRepo) Search(ctx context.Context, options *SearchOptions, page httputils.PageRequest) ([]*entity.Script, int64, error) {
	if options.Keyword != "" {
		// 暂时不支持排序等
		return u.SearchByEs(ctx, options, page)
	}
	// 无关键字从mysql数据库中获取
	list := make([]*entity.Script, 0)
	scriptTbName := (&entity.Script{}).TableName()
	find := db.Ctx(ctx).Model(&entity.Script{}).Where(scriptTbName+".status=?", consts.ACTIVE)
	if !options.Self {
		find.Where("public=? and unwell=?", entity.PublicScript, entity.Well)
	}
	switch options.Type {
	case 1: // 用户脚本
		find = find.Where("type=1")
	case 2: // 库
		find = find.Where("type=3")
	case 3: // 后台脚本
		options.Category = append(options.Category, 1)
	case 4: // 定时脚本
		options.Category = append(options.Category, 2)
	}
	if len(options.Category) != 0 {
		tabname := db.Default().NamingStrategy.TableName("script_category")
		find = find.Joins("left join "+tabname+" on "+tabname+".script_id="+scriptTbName+".id").
			Where(tabname+".category_id in ?", options.Category)
	}
	if options.Domain != "" {
		// 截取域名支持匹配子域名
		domains := strings.Split(utils.StringReverse(options.Domain), ".")
		tabname := db.Default().NamingStrategy.TableName("script_domain")
		find = find.Joins("left join "+tabname+" on "+tabname+".script_id="+scriptTbName+".id").
			Where(tabname+".status=?", consts.ACTIVE)
		if len(domains) <= 2 {
			find = find.Where(tabname+".domain=?", options.Domain)
		} else {
			exps := make([]clause.Expression, 0)
			for i := 1; i < len(domains); i++ {
				exps = append(exps, gorm.Expr(tabname+".domain_reverse like ?",
					strings.Join(domains[:i+1], ".")+"%"))
			}
			find = find.Where(clause.Or(exps...))
		}
	}

	switch options.Sort {
	case "today_download":
		tabname := db.Default().NamingStrategy.TableName("script_date_statistics")
		find = find.Joins(fmt.Sprintf("left join %s on %s.script_id=%s.id and %s.date=?", tabname, tabname, scriptTbName, tabname), time.Now().Format("2006-01-02")).
			Order(tabname + ".download desc,createtime desc")
	case "total_download":
		tabname := db.Default().NamingStrategy.TableName("script_statistics")
		find = find.Joins(fmt.Sprintf("left join %s on %s.script_id=%s.id", tabname, tabname, scriptTbName)).
			Order(tabname + ".download desc,createtime desc")
	case "score":
		tabname := db.Default().NamingStrategy.TableName("script_statistics")
		find = find.Joins(fmt.Sprintf("left join %s on %s.script_id=%s.id", tabname, tabname, scriptTbName)).
			Order(tabname + ".score desc,createtime desc")
	case "updatetime":
		find = find.Where(scriptTbName + ".updatetime>0").Order(scriptTbName + ".updatetime desc," + scriptTbName + ".createtime desc")
	default:
		find = find.Order(scriptTbName + ".createtime desc")
	}

	if options.UserID != 0 {
		find = find.Where("user_id=?", options.UserID)
	}
	var num int64
	if err := find.Count(&num).Error; err != nil {
		return nil, 0, err
	}
	find = find.Select(scriptTbName + ".*")
	find = find.Limit(page.GetLimit()).Offset(page.GetOffset())
	if err := find.Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, num, nil
}

// SearchByEs 通过elasticsearch搜索
func (u *scriptRepo) SearchByEs(ctx context.Context, options *SearchOptions, page httputils.PageRequest) ([]*entity.Script, int64, error) {
	script := &entity.ScriptSearch{}
	search := elasticsearch.Ctx(ctx).Search
	must := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":  options.Keyword,
				"fields": []string{"name", "description", "content"},
			},
		},
		{
			"match": map[string]interface{}{
				"status": consts.ACTIVE,
			},
		},
	}
	if !options.Self {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"public": entity.PublicScript,
			},
		}, map[string]interface{}{
			"match": map[string]interface{}{
				"unwell": entity.Well,
			},
		})
	}
	if options.UserID != 0 {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"user_id": options.UserID,
			},
		})
	}
	functionSearch := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}
	switch options.Sort {
	case "today_download":
		functionSearch["script_score"] = map[string]interface{}{
			"script": map[string]interface{}{
				"source": "Math.sqrt(10 + doc['today_download'].value)",
			},
		}
	case "score":
		functionSearch["script_score"] = map[string]interface{}{
			"script": map[string]interface{}{
				"source": "(doc['score'].value+10)/1.5",
			},
		}
	case "createtime":
		functionSearch["functions"] = []interface{}{
			map[string]interface{}{
				"gauss": map[string]interface{}{
					"createtime": map[string]interface{}{
						"origin": time.Now().Unix(),
						"offset": 1209600,
						"scale":  38880000,
						"decay":  0.33,
					},
				},
				"weight": 10,
			},
		}
	case "updatetime":
		functionSearch["functions"] = []interface{}{
			map[string]interface{}{
				"gauss": map[string]interface{}{
					"updatetime": map[string]interface{}{
						"origin": time.Now().Unix(),
						"offset": 1209600,
						"scale":  38880000,
						"decay":  0.33,
					},
				},
				"weight": 10,
			},
		}
	case "total_download":
		functionSearch["script_score"] = map[string]interface{}{
			"script": map[string]interface{}{
				"source": "_score * 200 + doc['today_download'].value",
			},
		}
	default:
		functionSearch["script_score"] = map[string]interface{}{
			"script": map[string]interface{}{
				"source": "_score * 200 + doc['today_download'].value",
			},
		}

	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"function_score": functionSearch,
		},
		"size": page.GetLimit(),
		"from": page.GetOffset(),
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
	ret := make([]*entity.Script, 0)
	for _, v := range m.Hits.Hits {
		source := v.Source
		ret = append(ret, &entity.Script{
			ID:     source.ID,
			PostID: 0,
			//PostID:        source.PostID,
			UserID:        source.UserID,
			Name:          source.Name,
			Description:   source.Description,
			Content:       source.Content,
			Type:          source.Type,
			Public:        source.Public,
			Unwell:        source.Unwell,
			SyncUrl:       "",
			ContentUrl:    "",
			DefinitionUrl: "",
			SyncMode:      0,
			Archive:       source.Archive,
			Status:        source.Status,
			Createtime:    source.Createtime,
			Updatetime:    source.Updatetime,
		})
	}
	return ret, m.Hits.Total.Value, nil
}

func (u *scriptRepo) FindSyncScript(ctx context.Context, page httputils.PageRequest) ([]*entity.Script, error) {
	var list []*entity.Script
	if err := db.Ctx(ctx).Where("sync_mode=? and status=? and sync_url<>''", entity.SyncModeAuto, consts.ACTIVE).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (u *scriptRepo) FindSyncPrefix(ctx context.Context, uid int64, prefix string) ([]*entity.Script, error) {
	var list []*entity.Script
	if err := db.Ctx(ctx).Where(
		"user_id=? and sync_mode=? and status=? and sync_url like ?",
		uid, entity.SyncModeAuto, consts.ACTIVE, prefix+"%",
	).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
