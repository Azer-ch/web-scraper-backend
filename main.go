package main

import (
	"github.com/Azer-ch/web-scraper/routes"
)

func main() {
	r := routes.SetupRouter()
	r.Run(":8080")
}
