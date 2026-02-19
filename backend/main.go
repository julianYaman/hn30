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
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"

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
	customUserAgent = "yamanlabs-hn/2.0 (+https://hn30.yamanlabs.com)"
)

var storyCache *Cache
var dbConn *sql.DB

var oneSignalConfig = onesignal.NewConfiguration()
var oneSignalApiClient = onesignal.NewAPIClient(oneSignalConfig)

func getTopStoryIDs() ([]int, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "fetch_top_story_ids",
		"url", fmt.Sprintf("%s/topstories.json", hnBaseURL),
	)
	start := time.Now()

	logger.Info("fetching top story ids",
		"event", "fetch_started",
	)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/topstories.json", hnBaseURL), nil)
	if err != nil {
		logger.Error("request creation failed",
			"event", "request_creation_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}
	req.Header.Set("User-Agent", customUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("http request failed",
			"event", "http_request_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("response received",
		"event", "response_received",
		"status_code", resp.StatusCode,
		"content_length", resp.ContentLength,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		logger.Error("json decode failed",
			"event", "json_decode_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}

	logger.Info("top story ids fetched successfully",
		"event", "fetch_completed",
		"total_ids", len(ids),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return ids, nil
}

func getStoryDetails(id int) (*types.Story, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "fetch_story_details",
		"story_id", id,
		"url", fmt.Sprintf("%s/item/%d.json", hnBaseURL, id),
	)
	start := time.Now()

	logger.Info("fetching story details",
		"event", "fetch_started",
	)

	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", hnBaseURL, id))
	if err != nil {
		logger.Error("http request failed",
			"event", "http_request_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("response received",
		"event", "response_received",
		"status_code", resp.StatusCode,
		"content_length", resp.ContentLength,
	)

	var story types.Story
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		logger.Error("json decode failed",
			"event", "json_decode_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}

	logger.Info("story details fetched successfully",
		"event", "fetch_completed",
		"story_title", story.Title,
		"story_url", story.URL,
		"story_score", story.Score,
		"story_by", story.By,
		"story_time", story.Time,
		"story_descendants", story.Descendants,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return &story, nil
}

func refreshCache() {

	ctx := context.Background()

	jobID := uuid.NewString()
	ctx = context.WithValue(ctx, "job_id", jobID)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("event_type", "cache_refresh", "job_id", jobID)
	cacheStart := time.Now()

	logger.Info(
		"cache refresh started",
		"event", "cache_refresh_started",
	)

	ids, err := getTopStoryIDs()
	if err != nil {
		logger.Error("cache_refresh_failed",
			"event", "cache_refresh_failed",
			"stage", "get_top_story_ids",
			"error", err,
		)
		return
	}

	topIDs := ids
	if len(topIDs) > 30 {
		topIDs = topIDs[:30]
	}

	storyCache.SetStoryIDs(topIDs)

	logger.Info("top_story_ids_fetched",
		"event", "top_story_ids_fetched",
		"total_ids", len(ids),
		"used_ids", len(topIDs),
	)

	for _, id := range topIDs {
		storyStart := time.Now()

		story, err := getStoryDetails(id)
		if err != nil {
			logger.Warn("story_processing_failed",
				"event", "story_processing_failed",
				"story_id", id,
				"stage", "get_story_details",
				"error", err,
			)
			continue
		}

		urlWasMissing := false
		if story.URL == "" {
			urlWasMissing = true
			story.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", id)
			// Fall through and attempt to fetch OG data for the HN item page
		}

		if existingStory, found := storyCache.Get(id); found && existingStory.URL == story.URL {
			existingStory.Score = story.Score
			existingStory.Descendants = story.Descendants

			utils.LogInfo("Story %d already in cache and URL unchanged, reusing OG data and updating stats", id)

			if err := db.UpsertStory(dbConn, *story); err != nil {
				logger.Error("story_upsert_failed",
					"event", "story_upsert_failed",
					"story_id", id,
					"cached", true,
					"error", err,
				)
			}

			if db.ShouldNotify(dbConn, *story) {
				go sendNotification(existingStory)
				db.MarkNotified(dbConn, story.ID)

				logger.Info("notification_sent",
					"event", "notification_sent",
					"story_id", id,
					"cached", true,
				)
			}

			storyCache.Set(id, existingStory)

			logger.Info("story_skipped_cached",
				"event", "story_skipped_cached",
				"story_id", id,
				"url", story.URL,
				"url_was_missing", urlWasMissing,
				"score", story.Score,
				"descendants", story.Descendants,
				"duration_ms", time.Since(storyStart).Milliseconds(),
			)

			continue
		}

		ogImage, ogDescription, err := getOGData(story.URL)
		if err != nil {
			logger.Warn("og_fetch_failed",
				"event", "og_fetch_failed",
				"story_id", id,
				"url", story.URL,
				"error", err,
			)
		}

		enrichedStory := EnrichedStory{
			Story:         *story,
			OGImage:       ogImage,
			OGDescription: ogDescription,
		}

		if err := db.UpsertStory(dbConn, *story); err != nil {
			logger.Error("story_upsert_failed",
				"event", "story_upsert_failed",
				"story_id", id,
				"cached", false,
				"error", err,
			)
		}

		notified := false
		if db.ShouldNotify(dbConn, *story) {
			go sendNotification(enrichedStory)
			db.MarkNotified(dbConn, story.ID)
			notified = true
		}

		storyCache.Set(id, enrichedStory)

		logger.Info("story_processed",
			"event", "story_processed",
			"story_id", id,
			"url", story.URL,
			"url_was_missing", urlWasMissing,
			"score", story.Score,
			"descendants", story.Descendants,
			"og_image_present", ogImage != "",
			"og_description_present", ogDescription != "",
			"notified", notified,
			"duration_ms", time.Since(storyStart).Milliseconds(),
		)

		time.Sleep(500 * time.Millisecond) // Rate limit scraping
	}

	storyCache.SetLastUpdated(time.Now())

	// Sync top stories to Turso for edge function OG image generation
	storiesToSync := make([]types.Story, 0, len(topIDs))
	for _, id := range topIDs {
		if s, found := storyCache.Get(id); found {
			storiesToSync = append(storiesToSync, s.Story)
		}
	}
	if err := db.SyncTopStories(storiesToSync); err != nil {
		logger.Warn("turso_sync_failed",
			"event", "turso_sync_failed",
			"error", err,
		)
	} else if len(storiesToSync) > 0 {
		logger.Info("turso_sync_completed",
			"event", "turso_sync_completed",
			"synced_count", len(storiesToSync),
		)
	}

	logger.Info("cache_refresh_completed",
		"event", "cache_refresh_completed",
		"duration_ms", time.Since(cacheStart).Milliseconds(),
		"stories_processed", len(topIDs),
	)
}

func sendNotification(story EnrichedStory) {

	ctx := context.Background()

	jobID := uuid.NewString()
	ctx = context.WithValue(ctx, "job_id", jobID)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "push_notification",
		"job_id", jobID,
		"story_id", story.ID,
		"story_title", story.Title,
		"story_url", story.URL,
		"story_score", story.Score,
	)
	start := time.Now()

	appId := os.Getenv("ONESIGNAL_APP_ID")
	restApiKey := os.Getenv("ONESIGNAL_KEY")

	if appId == "" || restApiKey == "" {
		logger.Error("onesignal credentials missing",
			"event", "credentials_missing",
			"app_id_present", appId != "",
			"api_key_present", restApiKey != "",
		)
		return
	}

	logger.Info("sending push notification",
		"event", "notification_send_started",
		"target_url", story.URL+"?ref=hn30",
	)

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

	notification.SetWebPushTopic("hn30_notifications-" + strconv.Itoa(story.ID)) // prevent overriding
	notification.SetPriority(10)                                                 // for iOS

	logger.Info("notification payload prepared",
		"event", "payload_prepared",
		"segments", []string{"Total Subscriptions"},
		"heading", "Top Story on Hacker News",
	)

	resp, httpResp, err := oneSignalApiClient.DefaultApi.CreateNotification(osAuthCtx).Notification(notification).Execute()

	if err != nil {
		logger.Error("notification send failed",
			"event", "notification_send_failed",
			"error", err,
			"http_status", func() int {
				if httpResp != nil {
					return httpResp.StatusCode
				}
				return 0
			}(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return
	}

	logger.Info("notification sent successfully",
		"event", "notification_sent",
		"notification_id", resp.GetId(),
		"external_id", resp.GetExternalId(),
		"http_status", httpResp.StatusCode,
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func startCacheRefresher() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "cache_refresher",
	)

	logger.Info("initializing cache refresher",
		"event", "cache_refresher_initialized",
		"refresh_interval", "5m",
	)

	storyCache = NewCache()
	go func() {
		logger.Info("starting initial cache refresh",
			"event", "initial_refresh_started",
		)
		refreshCache()

		ticker := time.NewTicker(5 * time.Minute)
		logger.Info("cache refresher running",
			"event", "refresher_running",
			"interval", "5m",
		)

		for range ticker.C {
			logger.Info("periodic cache refresh triggered",
				"event", "periodic_refresh_triggered",
			)
			refreshCache()
		}
	}()
}

func main() {
	log.SetFlags(0)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "server_lifecycle",
		"version", "1.0",
	)

	startTime := time.Now()

	logger.Info("application starting",
		"event", "app_start",
		"timestamp", startTime,
	)

	// Database initialization
	sqlitePath := os.Getenv("SQLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "./data/hn30.db"
	}
	logger.Info("initializing database",
		"event", "db_init_started",
		"db_path", sqlitePath,
	)
	dbConn = db.Open(sqlitePath)
	logger.Info("database initialized",
		"event", "db_init_completed",
		"db_path", sqlitePath,
	)

	// Turso initialization (optional - for edge function OG image generation)
	tursoURL := os.Getenv("TURSO_DATABASE_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
	if tursoURL != "" && tursoToken != "" {
		logger.Info("initializing turso connection",
			"event", "turso_init_started",
		)
		if err := db.OpenTurso(tursoURL, tursoToken); err != nil {
			logger.Warn("turso connection failed, OG image sync disabled",
				"event", "turso_init_failed",
				"error", err,
			)
		} else {
			logger.Info("turso connected",
				"event", "turso_init_completed",
			)
		}
	} else {
		logger.Info("turso not configured, OG image sync disabled",
			"event", "turso_skipped",
		)
	}

	// Cache initialization
	logger.Info("starting cache refresher",
		"event", "cache_init_started",
	)
	startCacheRefresher()
	logger.Info("cache refresher started",
		"event", "cache_init_completed",
	)

	// HTTP server setup
	server := &http.Server{Addr: ":8080"}

	http.Handle("/api/top", LoggingMiddleware(http.HandlerFunc(topStoriesHandler)))
	http.Handle("/api/summarize", LoggingMiddleware(rateLimitMiddleware(http.HandlerFunc(summarizeHandler))))

	logger.Info("http routes registered",
		"event", "routes_registered",
		"routes", []string{"/api/top", "/api/summarize"},
	)

	go func() {
		logger.Info("server starting",
			"event", "server_start",
			"addr", ":8080",
			"startup_duration_ms", time.Since(startTime).Milliseconds(),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server listen failed",
				"event", "server_listen_failed",
				"error", err,
				"addr", ":8080",
			)
			log.Fatalf("could not listen on port 8080 %v", err)
		}
	}()

	logger.Info("server ready",
		"event", "server_ready",
		"addr", ":8080",
		"total_startup_duration_ms", time.Since(startTime).Milliseconds(),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop

	shutdownStart := time.Now()
	logger.Info("shutdown signal received",
		"event", "shutdown_signal_received",
		"signal", sig.String(),
		"uptime_seconds", time.Since(startTime).Seconds(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("shutting down server",
		"event", "shutdown_started",
		"timeout", "5s",
	)

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed",
			"event", "shutdown_failed",
			"error", err,
			"duration_ms", time.Since(shutdownStart).Milliseconds(),
		)
		log.Fatalf("Server shutdown failed: %v", err)
	}

	if dbConn != nil {
		logger.Info("closing database connection",
			"event", "db_close_started",
		)
		if err := dbConn.Close(); err != nil {
			logger.Error("database close failed",
				"event", "db_close_failed",
				"error", err,
			)
		} else {
			logger.Info("database closed",
				"event", "db_close_completed",
			)
		}
	}

	// Close Turso connection
	if err := db.CloseTurso(); err != nil {
		logger.Error("turso close failed",
			"event", "turso_close_failed",
			"error", err,
		)
	}

	logger.Info("server stopped",
		"event", "shutdown_completed",
		"shutdown_duration_ms", time.Since(shutdownStart).Milliseconds(),
		"total_uptime_seconds", time.Since(startTime).Seconds(),
	)
}
