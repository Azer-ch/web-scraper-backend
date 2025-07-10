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

	cached, id, err := helpers.GetCachedResult(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cache error: " + err.Error()})
		return
	}
	if cached != nil {
		response := types.AnalyzeAPIResponse{
			ID:                id,
			HTMLVersion:       cached.HTMLVersion,
			Title:             cached.Title,
			Headings:          cached.Headings,
			InternalLinks:     cached.InternalLinks,
			ExternalLinks:     cached.ExternalLinks,
			InaccessibleLinks: cached.InaccessibleLinks,
			HasLoginForm:      cached.HasLoginForm,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	resp, err := helpers.AnalyzePage(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err = helpers.SetCachedResult(req.URL, resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := types.AnalyzeAPIResponse{
		ID:                id,
		HTMLVersion:       resp.HTMLVersion,
		Title:             resp.Title,
		Headings:          resp.Headings,
		InternalLinks:     resp.InternalLinks,
		ExternalLinks:     resp.ExternalLinks,
		InaccessibleLinks: resp.InaccessibleLinks,
		HasLoginForm:      resp.HasLoginForm,
	}
	c.JSON(http.StatusOK, response)
}
