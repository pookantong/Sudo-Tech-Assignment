package main

import (
	"log"

	"github.com/joho/godotenv"

	"cinema-booking-backend/internal/bootstrap"
	"cinema-booking-backend/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatalf("failed to load .env: %v", err)
    }
	
	cfg := config.Load()

	app := bootstrap.NewApp(cfg)

	log.Printf("Listening on :%s", cfg.Port)

	if err := app.Router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}