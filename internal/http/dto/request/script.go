package request

type Score struct {
	Score   int64  `form:"score" binding:"required,number,max=50,min=0" label:"分数" json:"score"`
	Message string `form:"message" binding:"required,max=200" label:"评论" json:"message"`
}

type CreateScript struct {
	Content     string `form:"content" binding:"required,max=102400" label:"脚本描述"`
	Code        string `form:"code" binding:"required,max=10485760" label:"脚本代码"`
	Name        string `form:"name" binding:"max=128" label:"库的名字"`
	Description string `form:"description" binding:"max=102400" label:"库的描述"`
	Definition  string `form:"definition" binding:"max=102400" label:"库的定义文件"`
	// 脚本类型：1 用户脚本 2 脚本调用库 3 订阅脚本
	Type int `form:"type" binding:"required" label:"脚本类型"`
	// 公开类型：1 公开 2 半公开
	Public int `form:"public" binding:"required" label:"公开类型"`
	Unwell int `form:"unwell" binding:"required" label:"不适内容"`
}
