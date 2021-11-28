package service

import "github.com/scriptscat/scriptlist/internal/domain/script/repository"

type ScriptWatch interface {
	// Watch 关注script
	Watch(script, user int64, level int) error
	Unwatch(script, user int64) error

	WatchList(script int64) (map[int64]int, error)
	IsWatch(script, user int64) (int, error)
}

type watch struct {
	watchRepo repository.ScriptWatch
}

func NewWatch(watchRepo repository.ScriptWatch) ScriptWatch {
	return &watch{watchRepo: watchRepo}
}

func (w *watch) Watch(script, user int64, level int) error {
	return w.watchRepo.Watch(script, user, level)
}

func (w *watch) Unwatch(script, user int64) error {
	return w.watchRepo.Unwatch(script, user)
}

func (w *watch) WatchList(script int64) (map[int64]int, error) {
	list, err := w.watchRepo.List(script)
	if err != nil {
		return nil, err
	}
	ret := make(map[int64]int, len(list))
	for _, v := range list {
		ret[v.UserId] = v.Level
	}
	return ret, nil
}

func (w *watch) IsWatch(script, user int64) (int, error) {
	return w.watchRepo.IsWatch(script, user)
}
