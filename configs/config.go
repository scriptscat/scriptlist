package configs

import (
	"context"

	"github.com/codfrm/cago/configs"
)

func Url() string {
	return configs.Default().String(context.Background(), "website.url")
}
