package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hn30/backend/db"
	"hn30/backend/types"
	"hn30/backend/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OneSignal/onesignal-go-api/v5"
)

type EnrichedStory struct {
	types.Story
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
var dbConn *sql.DB

var oneSignalConfig = onesignal.NewConfiguration()
var oneSignalApiClient = onesignal.NewAPIClient(oneSignalConfig)

func getTopStoryIDs() ([]int, error) {
	utils.LogComponent("HN_API", "Fetching top story IDs")
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
	utils.LogComponent("HN_API", "Successfully fetched %d top story IDs", len(ids))
	return ids, nil
}

func getStoryDetails(id int) (*types.Story, error) {
	utils.LogComponent("HN_API", "Fetching details for story %d", id)
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", hnBaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var story types.Story
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return nil, err
	}
	return &story, nil
}

func refreshCache() {
	utils.LogComponent("CACHE", "Starting cache refresh cycle")
	ids, err := getTopStoryIDs()
	if err != nil {
		utils.LogError("Failed to get top story IDs: %v", err)
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
			utils.LogWarn("Failed to get story details for ID %d: %v", id, err)
			continue
		}

		if story.URL == "" {
			// If the story has no external URL, use the Hacker News item page as the URL
			utils.LogInfo("Story %d has no URL; using HN item page as URL", id)
			story.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", id)
			// Fall through and attempt to fetch OG data for the HN item page
		}

		if existingStory, found := storyCache.Get(id); found && existingStory.URL == story.URL {
			utils.LogInfo("Story %d already in cache and URL unchanged, reusing OG data and updating stats", id)
			existingStory.Score = story.Score
			existingStory.Descendants = story.Descendants
			storyCache.Set(id, existingStory)
			continue
		}

		ogImage, ogDescription, err := getOGData(story.URL)
		if err != nil {
			utils.LogWarn("Failed to get OG data for URL %s: %v", story.URL, err)
		}

		enrichedStory := EnrichedStory{
			Story:         *story,
			OGImage:       ogImage,
			OGDescription: ogDescription,
		}

		dbErr := db.UpsertStory(dbConn, *story)
		if dbErr != nil {
			utils.LogError("Failed to upsert story %d into database: %v", id, dbErr)
		}

		if db.ShouldNotify(dbConn, *story) {
			go sendNotification(enrichedStory)
			db.MarkNotified(dbConn, story.ID)
		}

		storyCache.Set(id, enrichedStory)
		time.Sleep(500 * time.Millisecond) // Rate limit scraping
	}
	storyCache.SetLastUpdated(time.Now())
	utils.LogComponent("CACHE", "Cache refresh cycle complete")
}

func sendNotification(story EnrichedStory) {

	appId := os.Getenv("ONESIGNAL_APP_ID")
	restApiKey := os.Getenv("ONESIGNAL_KEY")

	utils.LogComponent("NOTIFICATION", "Sending notification for story %d", story.ID)

	osAuthCtx := context.WithValue(
		context.Background(),
		onesignal.RestApiKey,
		restApiKey)

	notification := *onesignal.NewNotification(appId)

	notification.SetUrl(story.URL + "?ref=hn30")

	content := onesignal.NewLanguageStringMap()
	content.SetEn(story.Title)
	notification.SetContents(*content)

	notification.SetIncludedSegments([]string{"Total Subscriptions"})

	headings := onesignal.NewLanguageStringMap()
	headings.SetEn("Top Story on Hacker News")
	notification.SetHeadings(*headings)

	// resp, r, err
	resp, _, err := oneSignalApiClient.DefaultApi.CreateNotification(osAuthCtx).Notification(notification).Execute()

	if err != nil {
		utils.LogError("Error: %v\n", err)
		return
	}

	fmt.Println("Response JSON:", resp)

	utils.LogComponent("NOTIFICATION", "Broadcast sent! ID: %v", resp.GetId())

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

	dbConn = db.Open(os.Getenv("SQLITE_PATH"))
	startCacheRefresher()

	server := &http.Server{Addr: ":8080"}

	http.Handle("/api/top", utils.LoggingMiddleware(http.HandlerFunc(topStoriesHandler)))
	http.Handle("/api/summarize", utils.LoggingMiddleware(rateLimitMiddleware(http.HandlerFunc(summarizeHandler))))

	go func() {
		utils.LogInfo("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on port 8080 %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	utils.LogInfo("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	utils.LogInfo("Server gracefully stopped")
}
