package request

type Score struct {
	Score   int64  `form:"score" binding:"required,number,max=50,min=0" label:"分数" json:"score"`
	Message string `form:"message" binding:"required,max=200" label:"评论" json:"message"`
}
