package configs

import (
	"context"

	"github.com/cago-frame/cago/configs"
)

var url string

func Url() string {
	if url != "" {
		return url
	}
	url = configs.Default().String(context.Background(), "website.url")
	if url == "" {
		url = "https://scriptcat.org"
	}
	return url
}
