package csrf

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

const CsrfSecret = "WdUiz9wR0WsufjgKhh1hrfApfXrXG854"

const Secret = "NQ3kDBBjRmBpRHSX3"

var _csrf = csrf.Protect([]byte(CsrfSecret))

func Token(c *gin.Context) string {
	return csrf.Token(c.Request)
}

func CsrfMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_csrf(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}
