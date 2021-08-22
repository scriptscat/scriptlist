package migrations

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/domain/script/repository"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

func DealMetaInfo() error {
	//创建后台脚本栏目
	categoryRepo := repository.NewCategory()
	bgCategory := &entity.ScriptCategoryList{
		Name:       "后台脚本",
		Createtime: time.Now().Unix(),
	}
	cronCategory := &entity.ScriptCategoryList{
		Name:       "定时脚本",
		Createtime: time.Now().Unix(),
	}
	if err := categoryRepo.Save(bgCategory); err != nil {
		return err
	}
	if err := categoryRepo.Save(cronCategory); err != nil {
		return err
	}

	for {
		list := make([]*entity.ScriptCode, 0)
		if err := db.Db.Model(&entity.ScriptCode{}).Where("meta_json is null").Limit(20).Scan(&list).Error; err != nil {
			return err
		}
		if len(list) == 0 {
			break
		}
		for _, v := range list {
			data := utils.ParseMetaToJson(v.Meta)
			jsonData, _ := json.Marshal(data)
			if err := db.Db.Model(&entity.ScriptCode{ID: v.ID}).Update("meta_json", string(jsonData)).Error; err != nil {
				return err
			}
			domains := make(map[string]struct{})
			if _, ok := data["background"]; ok {
				_ = categoryRepo.LinkCategory(v.ScriptId, bgCategory.ID)
			}
			if _, ok := data["crontab"]; ok {
				_ = categoryRepo.LinkCategory(v.ScriptId, bgCategory.ID)
				_ = categoryRepo.LinkCategory(v.ScriptId, cronCategory.ID)
			}
			for _, u := range data["match"] {
				domain := utils.ParseMetaDomain(u)
				if domain != "" {
					domains[domain] = struct{}{}
				} else {
					glog.Warningf("deal meta url info: %d %s", v.ID, u)
				}
			}
			for _, u := range data["include"] {
				domain := utils.ParseMetaDomain(u)
				if domain != "" {
					domains[domain] = struct{}{}
				} else {
					glog.Warningf("deal meta url info: %d %s", v.ID, u)
				}
			}
			for domain := range domains {
				db.Db.Save(&entity.ScriptDomain{
					Domain:        domain,
					DomainReverse: utils.StringReverse(domain),
					ScriptId:      v.ScriptId,
					ScriptCodeId:  v.ID,
					Createtime:    time.Now().Unix(),
				})
			}
		}
	}
	return nil
}
