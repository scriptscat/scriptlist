package service

import (
	"github.com/golang/glog"
	service2 "github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/statistics/service"
)

type Statistics interface {
	// Record 数据统计
	Record(scriptId, scriptCodeId, user int64, ip, ua, statisticsToken string, download string) error
}

type statistics struct {
	service.Statistics
	serviceSvc service2.Script
	queue      chan *recordParam
}

type recordParam struct {
	scriptId, scriptCodeId, user int64
	ip, ua, statisticsToken      string
	download                     string
}

func NewStatistics(statisSvc service.Statistics, serviceSvc service2.Script) Statistics {
	ret := &statistics{
		Statistics: statisSvc,
		serviceSvc: serviceSvc,
		queue:      make(chan *recordParam, 1000),
	}
	go ret.handlerQueue()
	return ret
}

func (s *statistics) handlerQueue() {
	for {
		record := <-s.queue
		if ok, err := s.Statistics.Record(record.scriptId, record.scriptCodeId, record.user, record.ip, record.ua, record.statisticsToken, record.download); err != nil {
			glog.Warningf("statis record error: %v", err)
		} else if !ok {
			continue
		}
		if record.download == service.DOWNLOAD_STATISTICS {
			if err := s.serviceSvc.Download(record.scriptId); err != nil {
				glog.Warningf("script statis record download: %v", err)
			}
		} else if record.download == service.UPDATE_STATISTICS {
			if err := s.serviceSvc.Update(record.scriptId); err != nil {
				glog.Warningf("script statis record update: %v", err)
			}
		}
	}
}

func (s *statistics) Record(scriptId, scriptCodeId, user int64, ip, ua, statisticsToken string, download string) error {
	s.queue <- &recordParam{
		scriptId:        scriptId,
		scriptCodeId:    scriptCodeId,
		user:            user,
		ip:              ip,
		ua:              ua,
		statisticsToken: statisticsToken,
		download:        download,
	}
	return nil
}
