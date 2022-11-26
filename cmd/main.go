package main

import (
	"log"

	"forum/internal/app"
	"forum/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
