package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var scraperClient = &http.Client{
	Timeout: 20 * time.Second, // Reduced from 30s to fail faster
}

func getOGData(storyUrl string) (string, string, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "scraper_operation",
		"operation", "get_og_data",
		"url", storyUrl,
	)
	start := time.Now()

	// Quick validation to fail fast
	if storyUrl == "" {
		logger.Error("empty url provided",
			"event", "validation_failed",
		)
		return "", "", fmt.Errorf("empty URL")
	}

	// Parse and validate URL before making request
	parsedURL, err := url.Parse(storyUrl)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		logger.Error("invalid url",
			"event", "validation_failed",
			"error", err,
		)
		return "", "", fmt.Errorf("invalid URL scheme")
	}

	logger.Info("scraping og data",
		"event", "scrape_started",
		"domain", parsedURL.Host,
	)

	// Create context with timeout for better control
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// Use an explicit request so we can set a proper User-Agent.
	req, err := http.NewRequestWithContext(ctx, "GET", storyUrl, nil)
	if err != nil {
		logger.Error("request creation failed",
			"event", "request_creation_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return "", "", err
	}

	// Be a polite scraper: identify the application and provide a URL.
	req.Header.Set("User-Agent", customUserAgent)

	res, err := scraperClient.Do(req)
	if err != nil {
		// Check if it's a timeout
		isTimeout := false
		if ctx.Err() == context.DeadlineExceeded {
			isTimeout = true
		}

		logger.Error("http request failed",
			"event", "http_request_failed",
			"error", err,
			"is_timeout", isTimeout,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return "", "", err
	}
	defer res.Body.Close()

	// Check for redirects to potentially problematic URLs
	if res.Request.URL.String() != storyUrl {
		logger.Info("redirect detected",
			"event", "redirect_detected",
			"original_url", storyUrl,
			"final_url", res.Request.URL.String(),
		)
	}

	// Be more lenient with status codes - some sites return OG data even on soft errors
	if res.StatusCode >= 400 {
		logger.Warn("error status code",
			"event", "error_status",
			"status_code", res.StatusCode,
			"status", res.Status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		// Still try to parse - some sites return OG data even on 404s
	}

	// Check content type - skip if not HTML-like
	contentType := res.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(strings.ToLower(contentType), "html") {
		logger.Info("non-html content type",
			"event", "content_type_skip",
			"content_type", contentType,
			"status_code", res.StatusCode,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return "", "", nil // Return empty but no error - not a failure
	}

	logger.Info("response received",
		"event", "response_received",
		"status_code", res.StatusCode,
		"content_type", contentType,
		"content_length", res.ContentLength,
	)

	// Limit response body size to prevent memory issues
	limitedReader := http.MaxBytesReader(nil, res.Body, 10*1024*1024) // 10MB max
	doc, err := goquery.NewDocumentFromReader(limitedReader)
	if err != nil {
		logger.Error("html parsing failed",
			"event", "html_parse_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return "", "", err
	}

	ogImage, _ := doc.Find("meta[property='og:image']").Attr("content")
	ogDescription, _ := doc.Find("meta[property='og:description']").Attr("content")

	// Resolve relative URLs to absolute
	if ogImage != "" {
		ogImage = resolveURL(storyUrl, ogImage, logger)
	}

	hasImage := ogImage != ""
	hasDescription := ogDescription != ""

	if hasImage || hasDescription {
		logger.Info("og data extracted successfully",
			"event", "scrape_completed",
			"has_image", hasImage,
			"has_description", hasDescription,
			"image_url", ogImage,
			"description", ogDescription,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	} else {
		logger.Info("no og data found",
			"event", "scrape_completed",
			"has_image", false,
			"has_description", false,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}

	return ogImage, ogDescription, nil
}

// resolveURL resolves relative URLs to absolute URLs
func resolveURL(baseURL, relativeURL string, logger *slog.Logger) string {
	if relativeURL == "" {
		return ""
	}

	// Try to parse as absolute URL first
	if parsed, err := url.Parse(relativeURL); err == nil && parsed.IsAbs() {
		return relativeURL
	}

	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		logger.Warn("failed to parse base url",
			"event", "url_parse_failed",
			"base_url", baseURL,
			"error", err,
		)
		return relativeURL
	}

	// Resolve relative URL
	resolved, err := base.Parse(relativeURL)
	if err != nil {
		logger.Warn("failed to resolve relative url",
			"event", "url_resolve_failed",
			"relative_url", relativeURL,
			"error", err,
		)
		return relativeURL
	}

	logger.Info("resolved relative url",
		"event", "url_resolved",
		"original", relativeURL,
		"resolved", resolved.String(),
	)

	return resolved.String()
}
