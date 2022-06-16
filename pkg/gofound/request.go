package gofound

type PutIndexRequest struct {
	ID       uint32      `json:"id"`
	Text     string      `json:"text"`
	Document interface{} `json:"document"`
}

type PutIndexBatchRequest []PutIndexRequest

type QueryOrder string

const (
	ORDER_DESC = "desc"
	ORDER_ASC  = "asc"
)

type QueryIndexRequest struct {
	Query     string     `json:"query"`
	Page      int        `json:"page,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Order     QueryOrder `json:"order,omitempty"`
	Highlight struct {
		PreTag  string `json:"preTag,omitempty"`
		PostTag string `json:"postTag,omitempty"`
	}
	ScoreExp string `json:"scoreExp,omitempty"`
}

type RemoveIndexRequest struct {
	ID uint32 `json:"id"`
}
