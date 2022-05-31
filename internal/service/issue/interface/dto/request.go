package dto

type UpdateLabels struct {
	Labels []string `json:"labels" binding:"required" label:"标签"`
}
