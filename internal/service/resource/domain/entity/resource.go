package entity

import "encoding/json"

type Resource struct {
	ID          string `json:"id"`
	Uid         int64  `json:"uid"`
	Comment     string `json:"comment"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	ContentType string `json:"content_type"`
	Createtime  int64  `json:"createtime"`
}

func (r *Resource) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}
