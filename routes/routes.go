package routes

import (
	"github.com/Azer-ch/web-scraper/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/analyze", handlers.AnalyzeHandler)
	return r
}
