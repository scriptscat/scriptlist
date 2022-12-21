package persistence

import (
	script2 "github.com/scriptscat/scriptlist/internal/repository/script"
)

func init() {
	script2.RegisterScript(new(script))
	script2.RegisterScriptCode(new(scriptCode))

	script2.RegisterScriptDomain(new(scriptDomain))
	script2.RegisterScriptCategory(new(scriptCategory))
	script2.RegisterScriptCategoryList(new(scriptCategoryList))

	script2.RegisterMigrate(new(migrate))
}
