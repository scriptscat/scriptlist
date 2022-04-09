package persistence

import (
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/repository"
	"gorm.io/gorm"
)

type follow struct {
	db *gorm.DB
}

func NewFollow(db *gorm.DB) repository.Follow {
	return &follow{
		db: db,
	}
}
func (f *follow) Find(uid, follow int64) (*entity.HomeFollow, error) {
	ret := &entity.HomeFollow{}
	if err := f.db.First(ret, "uid=? and followuid=?", uid, follow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (f *follow) List(uid int64, page *request.Pages) ([]*entity.HomeFollow, int64, error) {
	list := make([]*entity.HomeFollow, 0)
	find := f.db.Model(&entity.HomeFollow{}).Where("uid=?", uid)
	var count int64
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if page != request.AllPage {
		find = find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size())
	}
	if err := find.Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (f *follow) FollowerList(uid int64, page *request.Pages) ([]*entity.HomeFollow, int64, error) {
	list := make([]*entity.HomeFollow, 0)
	find := f.db.Model(&entity.HomeFollow{}).Where("followuid=?", uid)
	var count int64
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if page != request.AllPage {
		find = find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size())
	}
	if err := find.Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (f *follow) Save(homeFollow *entity.HomeFollow) error {
	return f.db.Create(homeFollow).Error
}

func (f *follow) Delete(uid, follow int64) error {
	return f.db.Delete(&entity.HomeFollow{}, "uid=? and followuid=?", uid, follow).Error
}
