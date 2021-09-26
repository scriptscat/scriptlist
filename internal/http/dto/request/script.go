package request

type Score struct {
	Score   int64  `form:"score" binding:"required,number,max=50,min=0" label:"分数" json:"score"`
	Message string `form:"message" binding:"required,max=200" label:"评论" json:"message"`
}

type CreateScript struct {
	Content     string `form:"content" binding:"required,max=102400" label:"脚本详细描述"`
	Code        string `form:"code" binding:"required,max=10485760" label:"脚本代码"`
	Name        string `form:"name" binding:"max=128" label:"库的名字"`
	Description string `form:"description" binding:"max=102400" label:"库的描述"`
	Definition  string `form:"definition" binding:"max=102400" label:"库的定义文件"`
	// 脚本类型：1 用户脚本 2 脚本调用库 3 订阅脚本
	Type int `form:"type" binding:"required" label:"脚本类型"`
	// 公开类型：1 公开 2 半公开
	Public    int    `form:"public" binding:"required" label:"公开类型"`
	Unwell    int    `form:"unwell" binding:"required" label:"不适内容"`
	Changelog string `form:"changelog" binding:"max=1024" label:"更新日志"`
}

type UpdateScript struct {
	Name        string `form:"name" binding:"max=128" label:"库的名字"`
	Description string `form:"description" binding:"max=102400" label:"库的描述"`
	// 公开类型：1 公开 2 半公开
	Public int `form:"public" binding:"required,number" label:"公开类型"`
	Unwell int `form:"unwell" binding:"required,number" label:"不适内容"`
	// 监听
	SyncUrl       string `form:"sync_url" binding:"omitempty,url,len=200" label:"代码同步url"`
	ContentUrl    string `form:"content_url" binding:"omitempty,url,len=200" label:"详细描述同步url"`
	DefinitionUrl string `form:"definition_url" binding:"omitempty,url,len=200" label:"定义文件同步url"`
	SyncMode      int    `form:"sync_mode" binding:"number" label:"同步模式"`
}

type UpdateScriptCode struct {
	Name        string `form:"name" binding:"max=128" label:"库的名字"`
	Description string `form:"description" binding:"max=102400" label:"库的描述"`
	Content     string `form:"content" binding:"required,max=102400" label:"脚本详细描述"`
	Code        string `form:"code" binding:"required,max=10485760" label:"脚本代码"`
	Definition  string `form:"definition" binding:"max=102400" label:"库的定义文件"`
	Changelog   string `form:"changelog" binding:"max=1024" label:"更新日志"`
	// 公开类型：1 公开 2 半公开
	Public int `form:"public" binding:"required,number" label:"公开类型"`
	Unwell int `form:"unwell" binding:"required,number" label:"不适内容"`
}
