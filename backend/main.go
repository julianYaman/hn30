package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Story struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Descendants int    `json:"descendants"`
}

type EnrichedStory struct {
	Story
	OGImage       string `json:"ogImage"`
	OGDescription string `json:"ogDescription"`
	Summary       string `json:"summary,omitempty"`
	ArticleText   string `json:"-"` // Don't send full text to client
	SummaryModel  string `json:"model,omitempty"`
}

const (
	hnBaseURL       = "https://hacker-news.firebaseio.com/v0"
	customUserAgent = "yamanlabs-hn/1.0 (+https://hn.yamanlabs.com)"
)

var storyCache *Cache

func getTopStoryIDs() ([]int, error) {
	LogComponent("HN_API", "Fetching top story IDs")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/topstories.json", hnBaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", customUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}
	LogComponent("HN_API", "Successfully fetched %d top story IDs", len(ids))
	return ids, nil
}

func getStoryDetails(id int) (*Story, error) {
	LogComponent("HN_API", "Fetching details for story %d", id)
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", hnBaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var story Story
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return nil, err
	}
	return &story, nil
}

func refreshCache() {
	LogComponent("CACHE", "Starting cache refresh cycle")
	ids, err := getTopStoryIDs()
	if err != nil {
		LogError("Failed to get top story IDs: %v", err)
		return
	}

	topIDs := ids
	if len(topIDs) > 30 {
		topIDs = topIDs[:30]
	}
	storyCache.SetStoryIDs(topIDs)

	for _, id := range topIDs {
		story, err := getStoryDetails(id)
		if err != nil {
			LogWarn("Failed to get story details for ID %d: %v", id, err)
			continue
		}

		if story.URL == "" {
			// If the story has no external URL, use the Hacker News item page as the URL
			LogInfo("Story %d has no URL; using HN item page as URL", id)
			story.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", id)
			// Fall through and attempt to fetch OG data for the HN item page
		}

		if existingStory, found := storyCache.Get(id); found && existingStory.URL == story.URL {
			LogInfo("Story %d already in cache and URL unchanged, reusing OG data and updating stats", id)
			existingStory.Score = story.Score
			existingStory.Descendants = story.Descendants
			storyCache.Set(id, existingStory)
			continue
		}

		ogImage, ogDescription, err := getOGData(story.URL)
		if err != nil {
			LogWarn("Failed to get OG data for URL %s: %v", story.URL, err)
		}

		enrichedStory := EnrichedStory{
			Story:         *story,
			OGImage:       ogImage,
			OGDescription: ogDescription,
		}
		storyCache.Set(id, enrichedStory)
		time.Sleep(500 * time.Millisecond) // Rate limit scraping
	}
	storyCache.SetLastUpdated(time.Now())
	LogComponent("CACHE", "Cache refresh cycle complete")
}

func startCacheRefresher() {
	storyCache = NewCache()
	go func() {
		refreshCache()
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			refreshCache()
		}
	}()
}

func main() {
	log.SetFlags(0)
	startCacheRefresher()

	server := &http.Server{Addr: ":8080"}

	http.Handle("/api/top", LoggingMiddleware(http.HandlerFunc(topStoriesHandler)))
	http.Handle("/api/summarize", LoggingMiddleware(rateLimitMiddleware(http.HandlerFunc(summarizeHandler))))
	http.Handle("/api/image-proxy", LoggingMiddleware(http.HandlerFunc(imageProxyHandler)))

	go func() {
		LogInfo("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on port 8080 %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	LogInfo("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	LogInfo("Server gracefully stopped")
}
