package gofound

type Response struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type PutIndexResponse Response

type RemoveIndexResponse Response

type QueryIndexResponse[T any] struct {
	Response
	Data *QueryIndexInfo[T] `json:"data"`
}

type QueryIndexInfo[T any] struct {
	Time      float64           `json:"time"`
	Total     int               `json:"total"`
	PageCount int               `json:"pageCount"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	Documents []DocumentItem[T] `json:"documents"`
	Words     []string          `json:"words"`
}

type DocumentItem[T any] struct {
	ID       uint32 `json:"id"`
	Text     string `json:"text"`
	Document T      `json:"document"`
	Score    int    `json:"score"`
}
