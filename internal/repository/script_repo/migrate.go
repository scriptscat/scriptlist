package script_repo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/pkg/logger"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_statistics_repo"
	"go.uber.org/zap"
)

type ScriptMigrateRepo interface {
	// SaveToEs 保存脚本数据到elasticsearch
	SaveToEs(ctx context.Context, s *entity.ScriptSearch) error
	// List 列出脚本数据
	List(ctx context.Context, start, size int) ([]*entity.Script, error)
	// Convert 转换为es储存的数据
	Convert(ctx context.Context, e *entity.Script) (*entity.ScriptSearch, error)
}

var defaultSearch ScriptMigrateRepo

func Migrate() ScriptMigrateRepo {
	return defaultSearch
}

func RegisterMigrate(i ScriptMigrateRepo) {
	defaultSearch = i
}

type migrateRepo struct {
}

func NewMigrateRepo() ScriptMigrateRepo {
	return &migrateRepo{}
}

func (m *migrateRepo) SaveToEs(ctx context.Context, s *entity.ScriptSearch) error {
	logger := logger.Ctx(ctx).With(zap.Int64("id", s.ID))
	r, err := s.Reader()
	if err != nil {
		return err
	}
	resp, err := elasticsearch.Ctx(ctx).Create(
		s.CollectionName(), strconv.FormatInt(s.ID, 10), r,
	)
	if err != nil {
		logger.Error("insert error", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		logger.Error("insert error", zap.ByteString("body", b), zap.Int("status", resp.StatusCode))
		return err
	}
	logger.Info("insert success")
	return nil
}

func (m *migrateRepo) List(ctx context.Context, start, size int) ([]*entity.Script, error) {
	list := make([]*entity.Script, 0, 20)
	if err := db.Ctx(ctx).Model(&entity.Script{}).
		Limit(size).Offset(start).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (m *migrateRepo) Convert(ctx context.Context, e *entity.Script) (*entity.ScriptSearch, error) {
	ret := &entity.ScriptSearch{
		ID:          e.ID,
		UserID:      e.UserID,
		Name:        e.Name,
		Description: e.Description,
		Content:     e.Content,
		Category:    nil,
		Domain:      nil,
		Public:      e.Public,
		Unwell:      e.Unwell,
		Createtime:  e.Createtime,
		Updatetime:  e.Updatetime,
	}
	code, err := ScriptCode().FindLatest(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, errors.New("code not found")
	}
	ret.Version = code.Version
	ret.Changelog = code.Changelog
	statistics, err := script_statistics_repo.ScriptStatistics().FindByScriptID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	if statistics != nil {
		ret.TotalDownload = statistics.Download
		ret.Score = float64(statistics.Score) / float64(statistics.ScoreCount)
	}
	dateStatistics, err := script_statistics_repo.ScriptDateStatistics().FindByScriptID(ctx, e.ID, time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	if dateStatistics != nil {
		ret.TodayDownload = dateStatistics.Download
	}
	list, err := ScriptCategory().List(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	ret.Category = make([]int64, 0, len(list))
	for _, v := range list {
		ret.Category = append(ret.Category, v.ID)
	}
	domain, err := Domain().List(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	ret.Domain = make([]string, 0)
	for _, v := range domain {
		ret.Domain = append(ret.Domain, v.Domain)
	}
	return ret, nil
}
