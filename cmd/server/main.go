package main

import (
	"log"
	"os"

	"go-basic/internal/config"
)

func main() {
	config.LoadEnv()

	app := config.NewGin()
	db := config.NewDatabase()
	validator := config.NewValidator()

	config.Bootstrap(&config.BootstrapConfig{
		App:       app,
		DB:        db,
		Validator: validator,
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
