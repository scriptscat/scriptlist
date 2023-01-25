package sender

import "github.com/scriptscat/scriptlist/internal/model/entity/user_entity"

type SendOptions struct {
	// From 发送者用户信息
	From *user_entity.User
	// Title 标题
	Title string
}
