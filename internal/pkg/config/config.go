package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mode    string
	Mysql   MySQL
	Redis   Redis
	Cache   Redis
	OAuth   OAuth `yaml:"oauth"`
	WebPort int   `yaml:"webPort"`
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

type OAuth struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
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
