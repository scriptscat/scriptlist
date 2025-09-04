package system

import "github.com/cago-frame/cago/server/mux"

// FeedbackRequest 用户反馈请求
type FeedbackRequest struct {
	mux.Meta `path:"/feedback" method:"POST"`
	// 遇到了错误或bug 不再需要脚本猫 缺少我需要的功能 找到了更好的替代品 其他原因
	Reason   string `json:"reason" binding:"required,oneof=bug unused feature better other"` // 反馈原因
	Content  string `json:"content" binding:"max=1000"`                                      // 反馈内容
	clientIp string // 客户端ip
}

func (f *FeedbackRequest) ClientIp() string {
	return f.clientIp
}

func (f *FeedbackRequest) SetClientIp(ip string) {
	f.clientIp = ip
}

type FeedbackResponse struct{}
