package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Mode    string
	MySQL   MySQL
	Redis   Redis
	WebPort int `yaml:"webPort"`
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type MySQL struct {
	Dsn    string
	Prefix string
}

var AppConfig Config

func Init(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("config read error: %v", err)

	}
	err = yaml.Unmarshal(file, &AppConfig)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}
	return nil
}
