package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-shiori/go-readability"
	"google.golang.org/genai"
)

var summarizerClient = &http.Client{
	Timeout: 10 * time.Second,
}

func generateSummary(articleText string) (string, error) {
	LogComponent("AI", "Generating summary using Google Gemini")
	if articleText == "" {
		return "", fmt.Errorf("cannot summarize empty text")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		LogComponent("AI", "Failed to create AI client: %v", err)
		return "", fmt.Errorf("failed to create AI client")
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-lite",
		genai.Text(`Summarize the following article in 2-5 concise sentences for a general audience, focusing on the main outcome. Format with line breaks. Omit links, references, and introductory phrases.

Article:
`+articleText),
		nil,
	)

	if err != nil {
		LogComponent("AI", "AI generation error: %v", err)
		return "", fmt.Errorf("AI generation error")
	}

	mockSummary := result.Text()
	LogComponent("AI", "Successfully generated summary")
	return mockSummary, nil
}

func extractArticleText(articleURL string) (string, error) {
	LogComponent("SCRAPER", "Extracting article text from %s", articleURL)

	// Create a new request so we can set headers
	req, err := http.NewRequest("GET", articleURL, nil)
	if err != nil {
		return "", err
	}

	// Set the custom user agent and other common browser headers
	req.Header.Set("User-Agent", customUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := summarizerClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: status %d", resp.StatusCode)
	}

	parsedURL, err := url.Parse(articleURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", err
	}

	LogComponent("SCRAPER", "Successfully extracted article text for %s", articleURL)
	return article.TextContent, nil
}
