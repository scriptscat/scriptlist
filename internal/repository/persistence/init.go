package persistence

import "github.com/scriptscat/scriptlist/internal/repository"

func init() {
	repository.RegisterUser(new(user))
}
