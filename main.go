package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xTaube/vr-controlled-robot-arm/server"
)


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// err = server.RunWebTransportServer(
		// os.Getenv("PORT"),
		// os.Getenv("CERT_PATH"),
		// os.Getenv("KEY_PATH"),
	// )
	err = server.RunWebSocketServer(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
