package request

type Score struct {
	Score   int64  `form:"score" binding:"required,number,max=50,min=0" label:"分数" json:"score"`
	Message string `form:"message" binding:"required,max=200" label:"评论" json:"message"`
}

type CreateScript struct {
	Description string `form:"description" binding:"required,max=102400" label:"脚本描述"`
	Code        string `form:"code" binding:"required,max=10485760" label:"脚本代码"`
	Definition  string `form:"definition" binding:"max=102400" label:"库的定义文件"`
	Type        int    `form:"type" binding:"required" label:"脚本类型"`
	Public      int    `form:"public" binding:"required" label:"公开类型"`
}
