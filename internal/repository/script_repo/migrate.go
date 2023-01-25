package script_repo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/pkg/logger"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/consts"
	"go.uber.org/zap"
)

// ScriptMigrateRepo 迁移到es或者其它数据库中
type ScriptMigrateRepo interface {
	// Save 保存脚本数据到elasticsearch
	Save(ctx context.Context, s *entity.ScriptSearch) error
	// List 列出脚本数据
	List(ctx context.Context, start, size int) ([]*entity.Script, error)
	// Convert 转换为es储存的数据
	Convert(ctx context.Context, e *entity.Script) (*entity.ScriptSearch, error)
	// Update 更新数据
	Update(ctx context.Context, s *entity.ScriptSearch) error
	// Delete 删除数据,但是是软删除
	Delete(ctx context.Context, id int64) error
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

func (m *migrateRepo) Save(ctx context.Context, s *entity.ScriptSearch) error {
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
		if resp.StatusCode == http.StatusConflict {
			// 更新
			return m.Update(ctx, s)
		}
		b, _ := io.ReadAll(resp.Body)
		logger.Error("insert error", zap.ByteString("body", b), zap.Int("status", resp.StatusCode))
		return fmt.Errorf("insert error: %d body: %s", resp.StatusCode, b)
	}
	logger.Info("insert success")
	return nil
}

func (m *migrateRepo) Update(ctx context.Context, s *entity.ScriptSearch) error {
	logger := logger.Ctx(ctx).With(zap.Int64("id", s.ID))
	r, err := s.Reader()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte("{\"doc\":"))
	_, _ = io.Copy(buf, r)
	buf.WriteString("}")
	resp, err := elasticsearch.Ctx(ctx).Update(
		s.CollectionName(), strconv.FormatInt(s.ID, 10), buf,
	)
	if err != nil {
		logger.Error("update error", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logger.Error("update error", zap.ByteString("body", b), zap.Int("status", resp.StatusCode))
		return fmt.Errorf("update error: %d body: %s", resp.StatusCode, b)
	}
	logger.Info("update success")
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
		Status:      e.Status,
		Createtime:  e.Createtime,
		Updatetime:  e.Updatetime,
	}
	code, err := ScriptCode().FindLatest(ctx, e.ID, false)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, errors.New("code not found")
	}
	ret.Version = code.Version
	ret.Changelog = code.Changelog
	statistics, err := ScriptStatistics().FindByScriptID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	if statistics != nil {
		ret.TotalDownload = statistics.Download
		if statistics.ScoreCount > 0 {
			ret.Score = float64(statistics.Score) / float64(statistics.ScoreCount)
		}
	}
	dateStatistics, err := ScriptDateStatistics().FindByScriptID(ctx, e.ID, time.Now())
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

func (m *migrateRepo) Delete(ctx context.Context, id int64) error {
	logger := logger.Ctx(ctx).With(zap.Int64("id", id))
	buf := bytes.NewBuffer([]byte(fmt.Sprintf("{\"doc\":{\"status\":%d}}", consts.DELETE)))
	resp, err := elasticsearch.Ctx(ctx).Update(
		(&entity.ScriptSearch{}).CollectionName(), strconv.FormatInt(id, 10), buf,
	)
	if err != nil {
		logger.Error("delete error", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logger.Error("delete error", zap.ByteString("body", b), zap.Int("status", resp.StatusCode))
		return fmt.Errorf("delete error: %d body: %s", resp.StatusCode, b)
	}
	logger.Info("delete success")
	return nil
}
