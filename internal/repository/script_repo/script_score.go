package script_repo

import (
	"context"
	"strconv"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/scriptscat/scriptlist/internal/model/entity"
)

type ScriptScoreRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptScore, error)
	Create(ctx context.Context, scriptScore *entity.ScriptScore) error
	Update(ctx context.Context, scriptScore *entity.ScriptScore) error
	Delete(ctx context.Context, id int64) error
	// ScoreList 获取评分列表
	ScoreList(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*entity.ScriptScore, int64, error)
	// FindByUser 查询该用户在该脚本下是否有过评分
	FindByUser(ctx context.Context, uid, scriptId int64) (*entity.ScriptScore, error)
	// LastScore 最新的评分
	LastScore(ctx context.Context, page httputils.PageRequest) ([]int64, error)
}

var defaultScriptScore ScriptScoreRepo

func ScriptScore() ScriptScoreRepo {
	return defaultScriptScore
}

func RegisterScriptScore(i ScriptScoreRepo) {
	defaultScriptScore = i
}
func NewScriptScore() ScriptScoreRepo {

	return &scriptScoreRepo{}
}

type scriptScoreRepo struct {
}

func (u *scriptScoreRepo) ScoreList(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*entity.ScriptScore, int64, error) {
	list := make([]*entity.ScriptScore, 0)
	find := db.Ctx(ctx).Model(&entity.ScriptScore{}).Where("script_id=? and state=?", scriptId, consts.ACTIVE).Order("createtime desc")
	var num int64
	if err := find.Count(&num).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Limit(page.GetLimit()).Offset(page.GetOffset()).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, num, nil
}

func (u *scriptScoreRepo) FindByUser(ctx context.Context, uid, scriptId int64) (*entity.ScriptScore, error) {
	ret := &entity.ScriptScore{}
	if err := db.Ctx(ctx).Where("user_id=? and script_id=?", uid, scriptId).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptScoreRepo) Find(ctx context.Context, id int64) (*entity.ScriptScore, error) {
	ret := &entity.ScriptScore{ID: id}
	if err := db.Ctx(ctx).Where("state=?", consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptScoreRepo) Create(ctx context.Context, scriptScore *entity.ScriptScore) error {
	if err := redis.Ctx(ctx).ZAdd("script:score:last", redis2.Z{
		Score:  float64(scriptScore.Createtime),
		Member: scriptScore.ScriptID,
	}).Err(); err != nil {
		return err
	}
	return db.Ctx(ctx).Create(scriptScore).Error
}

func (u *scriptScoreRepo) Update(ctx context.Context, scriptScore *entity.ScriptScore) error {
	return db.Ctx(ctx).Updates(scriptScore).Error
}

func (u *scriptScoreRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&entity.ScriptScore{ID: id}).Update("state", consts.DELETE).Error
}

func (u *scriptScoreRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*entity.ScriptScore, int64, error) {
	var list []*entity.ScriptScore
	var count int64
	find := db.Ctx(ctx).Model(&entity.ScriptScore{}).Where("state=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *scriptScoreRepo) LastScore(ctx context.Context, page httputils.PageRequest) ([]int64, error) {
	result, err := redis.Ctx(ctx).ZRevRangeByScore("script:score:last", &redis2.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  10,
	}).Result()
	if err != nil {
		return nil, err
	}
	list := make([]int64, 0)
	for _, v := range result {
		n, _ := strconv.ParseInt(v, 10, 64)
		list = append(list, n)
	}
	return list, nil
}
