package request

type Issue struct {
	Title   string   `form:"title" binding:"required,max=128" label:"标题"`
	Content string   `form:"content" binding:"max=10485760" label:"反馈内容"`
	Labels  []string `form:"labels" binding:"max=128" label:"标签"`
}

type Content struct {
	Content string `form:"content" binding:"required" label:"评论内容"`
}
