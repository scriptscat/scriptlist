package configs

import "github.com/codfrm/cago/configs"

func Url() string {
	return configs.Default().String("website.url")
}
