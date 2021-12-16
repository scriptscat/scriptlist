package migrations

import (
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/domain/script/repository"
	"github.com/scriptscat/scriptlist/internal/domain/script/service"
	repository2 "github.com/scriptscat/scriptlist/internal/domain/user/repository"
	service3 "github.com/scriptscat/scriptlist/internal/domain/user/service"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func T1638951531() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1638951531",
		Migrate: func(tx *gorm.DB) error {
			script := repository.NewScript()
			scriptWatch := repository.NewScriptWatch()
			scriptSvc := service.NewWatch(scriptWatch)
			userSvc := service3.NewUser(repository2.NewUser(), repository2.NewFollow())
			scripts, _, err := script.List(&repository.SearchList{Status: cnt.ACTIVE}, request.AllPage)
			if err != nil {
				logrus.Errorf("query script: %v", err)
				return err
			}
			email := make(map[string]*respond.User)
			for _, script := range scripts {
				if err := scriptSvc.Watch(script.ID, script.UserId, service.ScriptWatchLevelIssueComment); err != nil {
					logrus.Errorf("watch %v %v: %v", script.ID, script.UserId, err)
					continue
				}
				user, err := userSvc.SelfInfo(script.UserId)
				if err != nil {
					continue
				}
				email[user.Email] = user
			}
			logrus.Errorf("emails: %+v", email)
			msg := ""
			for _, v := range email {
				if v.Email == "" {
					continue
				}
				msg += v.Email + "\n"
			}
			if err := os.WriteFile("email.csv", []byte(msg), 0644); err != nil {
				logrus.Errorf("write email.csv %v: %v", msg, err)
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
