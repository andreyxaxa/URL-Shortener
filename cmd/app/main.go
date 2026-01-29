package main

import (
	"log"
	"os"

	"github.com/andreyxaxa/URL-Shortener/config"
	"github.com/andreyxaxa/URL-Shortener/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("config error: %s", err)
		}
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
