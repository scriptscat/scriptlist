package resource

import "github.com/codfrm/cago/server/mux"

// UploadImageRequest 上传图片
type UploadImageRequest struct {
	mux.Meta `path:"/resource/image" method:"POST"`
	Comment  string `form:"comment" binding:"required"`
	LinkID   int64  `form:"link_id" binding:"required"`
}

type UploadImageResponse struct {
	ID          string `json:"id"`
	LinkID      int64  `json:"link_id"`
	Comment     string `json:"comment"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Createtime  int64  `json:"createtime"`
}
