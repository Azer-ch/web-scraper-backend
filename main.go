package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type AnalyzeRequest struct {
	URL string `json:"url" binding:"required"`
}

type AnalyzeResponse struct {
	HTMLVersion       string         `json:"html_version"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internal_links"`
	ExternalLinks     int            `json:"external_links"`
	InaccessibleLinks int            `json:"inaccessible_links"`
	HasLoginForm      bool           `json:"has_login_form"`
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

	resp, err := analyzePage(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func analyzePage(pageURL string) (*AnalyzeResponse, error) {
	// Fetch the page
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// HTML Version
	htmlVersion := getHTMLVersion(pageURL)

	// Title
	title := doc.Find("title").Text()

	// Headings
	headings := map[string]int{}
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		headings[tag] = doc.Find(tag).Length()
	}

	// Links
	base, _ := url.Parse(pageURL)
	internalLinks := 0
	externalLinks := 0
	inaccessibleLinks := 0

	links := []string{}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		u, err := url.Parse(href)
		if err != nil || href == "" || strings.HasPrefix(href, "javascript:") {
			return
		}
		absURL := base.ResolveReference(u)
		if absURL.Host == base.Host {
			internalLinks++
		} else {
			externalLinks++
			links = append(links, absURL.String())
		}
	})

	// Check inaccessible links (4xx/5xx)
	for _, link := range links {
		status, err := getLinkStatus(link)
		if err != nil || (status >= 400 && status < 600) {
			inaccessibleLinks++
		}
	}

	// Login form detection
	hasLoginForm := false
	doc.Find("form").EachWithBreak(func(i int, s *goquery.Selection) bool {
		found := false
		s.Find("input[type='password']").Each(func(j int, _ *goquery.Selection) {
			found = true
		})
		if found {
			hasLoginForm = true
			return false
		}
		return true
	})

	return &AnalyzeResponse{
		HTMLVersion:       htmlVersion,
		Title:             title,
		Headings:          headings,
		InternalLinks:     internalLinks,
		ExternalLinks:     externalLinks,
		InaccessibleLinks: inaccessibleLinks,
		HasLoginForm:      hasLoginForm,
	}, nil
}

// Helper to get HTML version
func getHTMLVersion(pageURL string) string {
	resp, err := http.Get(pageURL)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return "Unknown"
		case html.DoctypeToken:
			token := tokenizer.Token()
			if strings.Contains(strings.ToLower(token.Data), "html") {
				return "HTML5"
			}
			return token.Data
		}
	}
}

// Helper to get link status
func getLinkStatus(link string) (int, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(link)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
