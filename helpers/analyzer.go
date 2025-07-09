package helpers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/Azer-ch/web-scraper/types"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func AnalyzePage(pageURL string) (*types.AnalyzeResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	htmlVersion := getHTMLVersion(pageURL)
	title := doc.Find("title").Text()
	headings := map[string]int{}
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		headings[tag] = doc.Find(tag).Length()
	}
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
	for _, link := range links {
		status, err := getLinkStatus(link)
		if err != nil || (status >= 400 && status < 600) {
			inaccessibleLinks++
		}
	}
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
	return &types.AnalyzeResponse{
		HTMLVersion:       htmlVersion,
		Title:             title,
		Headings:          headings,
		InternalLinks:     internalLinks,
		ExternalLinks:     externalLinks,
		InaccessibleLinks: inaccessibleLinks,
		HasLoginForm:      hasLoginForm,
	}, nil
}

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

func getLinkStatus(link string) (int, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(link)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
