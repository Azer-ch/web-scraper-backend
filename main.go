package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type AnalyzeRequest struct {
    URL string `json:"url" binding:"required"`
}

type AnalyzeResponse struct {
    HTMLVersion      string         `json:"html_version"`
    Title            string         `json:"title"`
    Headings         map[string]int `json:"headings"`
    InternalLinks    int            `json:"internal_links"`
    ExternalLinks    int            `json:"external_links"`
    InaccessibleLinks int           `json:"inaccessible_links"`
    HasLoginForm     bool           `json:"has_login_form"`
}

func main() {
    r := gin.Default()
    r.POST("/analyze", analyzeHandler)
    r.Run(":8080")
}

func analyzeHandler(c *gin.Context) {
    var req AnalyzeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp := AnalyzeResponse{
        HTMLVersion:      "HTML5",
        Title:            "Example Page",
        Headings:         map[string]int{"h1": 1, "h2": 2},
        InternalLinks:    5,
        ExternalLinks:    3,
        InaccessibleLinks: 1,
        HasLoginForm:     true,
    }
    c.JSON(http.StatusOK, resp)
}