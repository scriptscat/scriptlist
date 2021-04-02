package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Option func(opt *Options)

type Options struct {
	HandelHeader func(token *jwt.Token) error
}

func WithExpired(expired int64) func(opt *Options) {
	return func(opt *Options) {
		old := opt.HandelHeader
		opt.HandelHeader = func(token *jwt.Token) error {
			if old != nil {
				if err := old(token); err != nil {
					return err
				}
			}
			if t, ok := token.Header["time"]; ok {
				if i, ok := t.(float64); !(ok && int64(i)+expired > time.Now().Unix()) {
					return fmt.Errorf("token failure")
				}
			} else {
				return fmt.Errorf("token failure")
			}
			return nil
		}
	}
}
