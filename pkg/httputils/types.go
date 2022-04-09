package httputils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type List struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func QueryInt64(c *gin.Context, key string) int64 {
	i, _ := strconv.ParseInt(c.Query(key), 10, 64)
	return i
}

func BindMap(c *gin.Context) (map[string]string, error) {
	data := make(map[string]string)
	if err := c.ShouldBind(data); err != nil {
		return nil, err
	}
	return data, nil
}
