package application

import (
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
)

type ScriptWatchLevel int

const (
	ScriptWatchLevelNone ScriptWatchLevel = iota
	ScriptWatchLevelVersion
	ScriptWatchLevelIssue
	ScriptWatchLevelIssueComment
)

type ScriptWatch interface {
	// Watch 关注script
	Watch(script, user int64, level ScriptWatchLevel) error
	Unwatch(script, user int64) error

	WatchList(script int64) (map[int64]ScriptWatchLevel, error)
	IsWatch(script, user int64) (ScriptWatchLevel, error)
}

type watch struct {
	watchRepo repository.ScriptWatch
}

func NewWatch(watchRepo repository.ScriptWatch) ScriptWatch {
	return &watch{watchRepo: watchRepo}
}

func (w *watch) Watch(script, user int64, level ScriptWatchLevel) error {
	return w.watchRepo.Watch(script, user, int(level))
}

func (w *watch) Unwatch(script, user int64) error {
	return w.watchRepo.Unwatch(script, user)
}

func (w *watch) WatchList(script int64) (map[int64]ScriptWatchLevel, error) {
	list, err := w.watchRepo.List(script)
	if err != nil {
		return nil, err
	}
	ret := make(map[int64]ScriptWatchLevel, len(list))
	for _, v := range list {
		ret[v.UserId] = ScriptWatchLevel(v.Level)
	}
	return ret, nil
}

func (w *watch) IsWatch(script, user int64) (ScriptWatchLevel, error) {
	ret, err := w.watchRepo.IsWatch(script, user)
	return ScriptWatchLevel(ret), err
}
