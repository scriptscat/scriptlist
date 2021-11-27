package main

import (
	"flag"
	"log"

	"github.com/scriptscat/scriptlist/internal/app"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
)

func main() {
	cfg := "config.yaml"
	flag.StringVar(&cfg, "config", cfg, "配置文件")
	flag.Parse()
	if err := config.Init(cfg); err != nil {
		log.Fatal("config error: ", err)
	}
	app.Run()
}
