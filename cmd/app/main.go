package main

import (
	"log"

	"github.com/scriptscat/scriptweb/internal/app"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
)

func main() {
	if err := config.Init("config.yaml"); err != nil {
		log.Fatal("config error: ", err)
	}
	app.Run()
}
