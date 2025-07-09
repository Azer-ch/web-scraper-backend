package handlers

import (
	"net/http"
	"github.com/Azer-ch/web-scraper/helpers"
	"github.com/Azer-ch/web-scraper/types"
	"github.com/gin-gonic/gin"
)

func AnalyzeHandler(c *gin.Context) {
	var req types.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cached, err := helpers.GetCachedResult(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cache error: " + err.Error()})
		return
	}
	if cached != nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	resp, err := helpers.AnalyzePage(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = helpers.SetCachedResult(req.URL, resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
