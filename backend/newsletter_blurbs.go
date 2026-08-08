package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hn30/backend/db"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var blockedBlurbHosts = map[string]struct{}{
	"x.com": {}, "twitter.com": {},
	"instagram.com": {}, "facebook.com": {}, "tiktok.com": {},
}

// Longer timeout for extract + OpenRouter during digest generation.
var digestClient = &http.Client{
	Timeout: 45 * time.Second,
}

func digestModel() string {
	if m := strings.TrimSpace(os.Getenv("OPENROUTER_DIGEST_MODEL")); m != "" {
		return m
	}
	return "@preset/hn30-digest-blurb"
}

func isBlockedBlurbURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return true
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	_, blocked := blockedBlurbHosts[host]
	return blocked
}

// ensureStoryBlurb returns a blurb (may be empty). Failures never abort the digest.
// Permanent skips are only stored for blocked hosts. Transient extract/AI failures
// are not persisted, so later runs can retry.
func ensureStoryBlurb(hnID int, articleURL string, logger *slog.Logger) string {
	existing, err := db.GetStoryBlurb(dbConn, hnID)
	if err != nil {
		logger.Warn("blurb lookup failed", "hn_id", hnID, "error", err)
	} else if existing != nil && existing.Status == db.BlurbStatusOK && strings.TrimSpace(existing.Blurb) != "" {
		return existing.Blurb
	} else if existing != nil && existing.Status == db.BlurbStatusSkipped {
		// Only blocked hosts are stored as skipped permanently.
		return ""
	}

	if articleURL == "" || isBlockedBlurbURL(articleURL) {
		_ = db.UpsertStoryBlurb(dbConn, hnID, "", db.BlurbStatusSkipped, "")
		return ""
	}

	articleText, err := extractArticleTextWithClient(articleURL, digestClient)
	if err != nil || strings.TrimSpace(articleText) == "" {
		logger.Info("blurb skipped: extract failed",
			"hn_id", hnID,
			"url", articleURL,
			"error", errString(err),
		)
		return ""
	}

	if len(articleText) > 12000 {
		articleText = articleText[:12000]
	}

	prompt := "Write a plain-text preview of this article for a daily email digest. Exactly 1 or 2 short sentences. No title, no bullet points, no preamble.\n\n" + articleText
	summary, err := generateDigestBlurb(prompt)
	if err != nil || strings.TrimSpace(summary.Summary) == "" {
		logger.Info("blurb skipped: ai failed",
			"hn_id", hnID,
			"url", articleURL,
			"error", errString(err),
		)
		return ""
	}

	blurb := strings.TrimSpace(summary.Summary)
	blurb = strings.Trim(blurb, `"'`)
	if err := db.UpsertStoryBlurb(dbConn, hnID, blurb, db.BlurbStatusOK, summary.Model); err != nil {
		logger.Warn("blurb persist failed", "hn_id", hnID, "error", err)
	}
	return blurb
}

func generateDigestBlurb(prompt string) (SummaryResponse, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return SummaryResponse{}, fmt.Errorf("OPENROUTER_API_KEY not set")
	}

	model := digestModel()
	body := OpenRouterRequest{
		Model: model,
		Messages: []OpenRouterMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return SummaryResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return SummaryResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "https://hn30.eu")
	req.Header.Set("X-Title", "hn30-digest")

	res, err := digestClient.Do(req)
	if err != nil {
		return SummaryResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return SummaryResponse{}, fmt.Errorf("openrouter model=%s status %d: %s", model, res.StatusCode, truncate(string(b), 300))
	}

	var openRouterResp OpenRouterResponse
	if err := json.NewDecoder(res.Body).Decode(&openRouterResp); err != nil {
		return SummaryResponse{}, err
	}
	if len(openRouterResp.Choices) == 0 {
		return SummaryResponse{}, fmt.Errorf("no choices in response")
	}
	return SummaryResponse{
		Summary: openRouterResp.Choices[0].Message.Content,
		Model:   openRouterResp.Model,
	}, nil
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
