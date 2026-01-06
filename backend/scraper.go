package main

import (
	"fmt"
	"hn30/backend/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var scraperClient = &http.Client{
	Timeout: 30 * time.Second,
}

func getOGData(storyUrl string) (string, string, error) {
	utils.LogComponent("SCRAPER", "Scraping %s for OG data", storyUrl)
	// Use an explicit request so we can set a proper User-Agent.
	req, err := http.NewRequest("GET", storyUrl, nil)
	if err != nil {
		utils.LogWarn("Failed to create request for %s: %v", storyUrl, err)
		return "", "", err
	}

	// Be a polite scraper: identify the application and provide a URL.
	// Change the string below if you want a different identifier or contact URL.
	req.Header.Set("User-Agent", customUserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogWarn("Failed to fetch URL %s: %v", storyUrl, err)
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		utils.LogWarn("Scraping %s returned non-200 status: %d", storyUrl, res.StatusCode)
		return "", "", fmt.Errorf("bad status: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	ogImage, _ := doc.Find("meta[property='og:image']").Attr("content")
	ogDescription, _ := doc.Find("meta[property='og:description']").Attr("content")

	// If ogImage is a relative URL, resolve it to absolute
	if ogImage != "" && ogImage[0] == '/' {
		// Get origin from the input URL
		if parsed, err := url.Parse(storyUrl); err == nil {
			origin := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
			ogImage = origin + ogImage
		}
	}

	if ogImage != "" || ogDescription != "" {
		utils.LogComponent("SCRAPER", "Successfully scraped OG data from %s", storyUrl)
	} else {
		utils.LogComponent("SCRAPER", "No OG data found on %s", storyUrl)
	}

	return ogImage, ogDescription, nil
}
