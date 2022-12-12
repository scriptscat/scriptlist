package persistence

import (
	"github.com/scriptscat/scriptlist/internal/repository"
	script2 "github.com/scriptscat/scriptlist/internal/repository/script"
)

func init() {
	repository.RegisterUser(new(user))

	script2.RegisterScript(new(script))
	script2.RegisterScriptCode(new(scriptCode))

	script2.RegisterScriptDomain(new(scriptDomain))
	script2.RegisterScriptCategory(new(scriptCategory))
	script2.RegisterScriptCategoryList(new(scriptCategoryList))
}
