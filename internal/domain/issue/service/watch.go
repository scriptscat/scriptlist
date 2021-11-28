package service

import "github.com/scriptscat/scriptlist/internal/domain/issue/repository"

type Watch interface {
	// Watch 关注issue
	Watch(issue, user int64) error
	Unwatch(issue, user int64) error

	WatchList(issue int64) ([]int64, error)
	IsWatch(issue, user int64) (int, error)
}

type watch struct {
	watchRepo repository.IssueWatch
}

func NewWatch(watchRepo repository.IssueWatch) Watch {
	return &watch{watchRepo: watchRepo}
}

func (w *watch) Watch(issue, user int64) error {
	return w.watchRepo.Watch(issue, user)
}

func (w *watch) Unwatch(issue, user int64) error {
	return w.watchRepo.Unwatch(issue, user)
}

func (w *watch) WatchList(issue int64) ([]int64, error) {
	list, err := w.watchRepo.List(issue)
	if err != nil {
		return nil, err
	}
	ret := make([]int64, len(list))
	for k, v := range list {
		ret[k] = v.UserId
	}
	return ret, nil
}

func (w *watch) IsWatch(issue, user int64) (int, error) {
	return w.watchRepo.IsWatch(issue, user)
}
