package script_svc

import (
	"context"
	"strings"
	"time"

	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type CategorySvc interface {
	// CategoryList 获取脚本分类列表
	CategoryList(ctx context.Context, req *api.CategoryListRequest) (*api.CategoryListResponse, error)
	// LinkScriptCategory 关联脚本分类
	LinkScriptCategory(ctx context.Context, scriptId, categoryId int64) error
	// LinkScriptTag 关联脚本标签
	LinkScriptTag(ctx context.Context, scriptId int64, tags []string) error
}

type categorySvc struct {
}

var defaultCategory = &categorySvc{}

func Category() CategorySvc {
	return defaultCategory
}

// CategoryList 获取脚本分类列表
func (c *categorySvc) CategoryList(ctx context.Context, req *api.CategoryListRequest) (*api.CategoryListResponse, error) {
	list, err := script_repo.NewScriptCategoryListRepo().FindByNamePrefixAndType(ctx, req.Prefix, req.Type)
	if err != nil {
		return nil, err
	}
	response := &api.CategoryListResponse{
		Categories: make([]*api.CategoryListItem, 0, len(list)),
	}
	for _, item := range list {
		response.Categories = append(response.Categories, &api.CategoryListItem{
			ID:         item.ID,
			Name:       item.Name,
			Num:        item.Num,
			Sort:       item.Sort,
			Type:       item.Type,
			Createtime: item.Createtime,
			Updatetime: item.Updatetime,
		})
	}
	return response, nil
}

// LinkScriptCategory 关联脚本分类
func (c *categorySvc) LinkScriptCategory(ctx context.Context, scriptId, categoryId int64) error {
	if categoryId == 0 {
		// 删除原来的分类
		list, err := script_repo.NewScriptCategoryRepo().FindByScriptId(ctx, scriptId, script_entity.ScriptCategoryTypeCategory)
		if err != nil {
			return err
		}
		for _, item := range list {
			if err := script_repo.NewScriptCategoryRepo().Delete(ctx, item); err != nil {
				return err
			}
		}
		return nil
	}
	// 检查分类是否存在
	category, err := script_repo.NewScriptCategoryListRepo().Find(ctx, categoryId)
	if err != nil {
		return err
	}
	if category == nil || category.Type != script_entity.ScriptCategoryTypeCategory {
		return i18n.NewError(ctx, code.ScriptCategoryNotFound)
	}
	list, err := script_repo.NewScriptCategoryRepo().FindByScriptId(ctx, scriptId, script_entity.ScriptCategoryTypeCategory)
	if err != nil {
		return err
	}
	if len(list) != 0 {
		// 检查是否与脚本已关联
		if list[0].CategoryID != category.ID {
			return nil
		}
		// 删除关联
		if err := script_repo.NewScriptCategoryRepo().Delete(ctx, list[0]); err != nil {
			return err
		}
	}
	// 进行关联
	if err := script_repo.NewScriptCategoryRepo().Create(ctx, &script_entity.ScriptCategory{
		CategoryID: categoryId,
		ScriptID:   scriptId,
		Createtime: time.Now().Unix(),
	}); err != nil {
		return err
	}
	return nil
}

// LinkScriptTag 关联脚本标签
func (c *categorySvc) LinkScriptTag(ctx context.Context, scriptId int64, tags []string) error {
	if len(tags) == 0 {
		// 删除原来的标签
		list, err := script_repo.NewScriptCategoryRepo().FindByScriptId(ctx, scriptId, script_entity.ScriptCategoryTypeTag)
		if err != nil {
			return err
		}
		for _, item := range list {
			if err := script_repo.NewScriptCategoryRepo().Delete(ctx, item); err != nil {
				return err
			}
		}
		return nil
	}
	// 重新处理tag，根据,空格分隔并去重
	tagsMap := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSplit := strings.Split(tag, ", ")
		for _, t := range tagSplit {
			t = strings.TrimSpace(t)
			if t != "" {
				tagsMap[t] = struct{}{}
			}
		}
	}
	tags = make([]string, 0, len(tagsMap))
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	// 检查标签是否存在
	categoryList := make([]*script_entity.ScriptCategoryList, 0)
	for _, tag := range tags {
		categoryListItem, err := script_repo.NewScriptCategoryListRepo().FindByNameAndType(ctx, tag, script_entity.ScriptCategoryTypeTag)
		if err != nil {
			return err
		}
		if categoryListItem == nil {
			// 如果标签不存在，则创建
			categoryListItem = &script_entity.ScriptCategoryList{
				Name:       tag,
				Type:       script_entity.ScriptCategoryTypeTag,
				Createtime: time.Now().Unix(),
			}
			if err := script_repo.NewScriptCategoryListRepo().Create(ctx, categoryListItem); err != nil {
				return err
			}
		}
		categoryList = append(categoryList, categoryListItem)
	}
	if len(categoryList) == 0 {
		return i18n.NewError(ctx, code.ScriptCategoryNotFound)
	}
	// 获取脚本原来的标签列表
	list, err := script_repo.NewScriptCategoryRepo().FindByScriptId(ctx, scriptId, script_entity.ScriptCategoryTypeTag)
	if err != nil {
		return err
	}
	// 对比原来的标签列表和新的标签列表，删除原来没有的标签
	existMap := make(map[int64]struct{}, len(categoryList))
	for _, item := range categoryList {
		existMap[item.ID] = struct{}{}
	}
	for _, item := range list {
		if _, ok := existMap[item.CategoryID]; !ok {
			// 删除原来没有的标签
			if err := script_repo.NewScriptCategoryRepo().Delete(ctx, item); err != nil {
				return err
			}
		} else {
			// 如果标签存在于原来的列表中，则从新的列表中删除，避免重复添加
			delete(existMap, item.CategoryID)
		}
	}
	// 将新的标签列表添加到脚本中
	for categoryId := range existMap {
		// 添加新的标签
		if err := script_repo.NewScriptCategoryRepo().Create(ctx, &script_entity.ScriptCategory{
			CategoryID: categoryId,
			ScriptID:   scriptId,
			Createtime: time.Now().Unix(),
		}); err != nil {
			return err
		}
	}
	return nil
}
