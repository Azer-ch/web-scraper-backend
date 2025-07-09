package main

import (
	"github.com/Azer-ch/web-scraper/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	r := routes.SetupRouter()
	r.Run(":8080")
}
