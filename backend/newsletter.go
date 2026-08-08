package main

import (
	"fmt"
	"hn30/backend/db"
	"log/slog"
	"os"
	"strings"
	"time"
)

func newsletterEnabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("NEWSLETTER_ENABLED")))
	return v == "1" || v == "true" || v == "yes"
}

func berlinNow() (time.Time, error) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

func startNewsletterScheduler() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "newsletter_scheduler",
	)

	if !newsletterEnabled() {
		logger.Info("newsletter scheduler disabled",
			"event", "newsletter_disabled",
			"hint", "set NEWSLETTER_ENABLED=true to enable",
		)
		return
	}

	logger.Info("newsletter scheduler starting",
		"event", "newsletter_scheduler_started",
		"send_time", "07:30 Europe/Berlin",
	)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		var lastAttemptDate string

		for range ticker.C {
			now, err := berlinNow()
			if err != nil {
				logger.Error("berlin timezone load failed", "error", err)
				continue
			}

			digestDate := now.Format("2006-01-02")
			if now.Hour() != 7 || now.Minute() < 30 {
				continue
			}
			if lastAttemptDate == digestDate {
				continue
			}

			lastAttemptDate = digestDate
			logger.Info("newsletter window reached",
				"event", "newsletter_window",
				"digest_date", digestDate,
				"local_time", now.Format(time.RFC3339),
			)
			if err := runDailyDigest(digestDate, true); err != nil {
				logger.Error("newsletter job failed",
					"event", "newsletter_job_failed",
					"digest_date", digestDate,
					"error", err,
				)
				// Allow another attempt later the same morning.
				lastAttemptDate = ""
			}
		}
	}()
}

// runDailyDigest builds today's snapshot, ensures blurbs, and optionally sends via Mailjet.
func runDailyDigest(digestDate string, send bool) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "newsletter_job",
		"digest_date", digestDate,
		"send", send,
	)

	claimed, err := db.ClaimDigest(dbConn, digestDate)
	if err != nil {
		return fmt.Errorf("claim digest: %w", err)
	}
	if !claimed {
		logger.Info("digest already sent or not claimable",
			"event", "newsletter_skipped",
		)
		return nil
	}

	stories := storyCache.GetAll()
	if len(stories) == 0 {
		err := fmt.Errorf("no stories in cache")
		_ = db.MarkDigestFailed(dbConn, digestDate, err)
		return err
	}
	if len(stories) > 30 {
		stories = stories[:30]
	}

	items := make([]db.DigestItem, 0, len(stories))
	for i, story := range stories {
		blurb := ensureStoryBlurb(story.ID, story.URL, logger)
		items = append(items, db.DigestItem{
			DigestDate: digestDate,
			Rank:       i + 1,
			HNID:       story.ID,
			Title:      story.Title,
			URL:        story.URL,
			Score:      story.Score,
			Blurb:      blurb,
		})
	}

	if err := db.ReplaceDigestItems(dbConn, digestDate, items); err != nil {
		_ = db.MarkDigestFailed(dbConn, digestDate, err)
		return fmt.Errorf("persist digest items: %w", err)
	}

	emailStories := digestItemsToEmailStories(items)
	htmlBody := renderDigestHTML(digestDate, emailStories, "")
	textBody := renderDigestText(digestDate, emailStories)
	subject := "hn30 daily dispatch · " + formatDigestDateLabel(digestDate)
	title := "hn30-" + digestDate

	if !send {
		logger.Info("digest built without send",
			"event", "newsletter_built",
			"story_count", len(items),
		)
		return nil
	}

	client, err := newMailjetClientFromEnv()
	if err != nil {
		_ = db.MarkDigestFailed(dbConn, digestDate, err)
		return err
	}

	draftID, err := client.CreateAndSendCampaign(title, subject, htmlBody, textBody)
	if err != nil {
		_ = db.MarkDigestFailed(dbConn, digestDate, err)
		return err
	}

	if err := db.MarkDigestSent(dbConn, digestDate, draftID); err != nil {
		return fmt.Errorf("mark sent: %w", err)
	}

	logger.Info("newsletter sent",
		"event", "newsletter_sent",
		"mailjet_draft_id", draftID,
		"story_count", len(items),
	)
	return nil
}

// buildDigestHTML returns the email HTML for today's digest.
// Prefers a frozen morning snapshot when present; otherwise live top 30 + stored blurbs.
// Never calls OpenRouter.
func buildDigestHTML() (string, error) {
	now, err := berlinNow()
	if err != nil {
		return "", err
	}
	digestDate := now.Format("2006-01-02")

	items, err := db.GetDigestItems(dbConn, digestDate)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		stories := storyCache.GetAll()
		if len(stories) == 0 {
			return "", fmt.Errorf("no stories available")
		}
		if len(stories) > 30 {
			stories = stories[:30]
		}
		items = make([]db.DigestItem, 0, len(stories))
		for i, story := range stories {
			blurb := ""
			if existing, err := db.GetStoryBlurb(dbConn, story.ID); err == nil && existing != nil && existing.Status == db.BlurbStatusOK {
				blurb = existing.Blurb
			}
			items = append(items, db.DigestItem{
				DigestDate: digestDate,
				Rank:       i + 1,
				HNID:       story.ID,
				Title:      story.Title,
				URL:        story.URL,
				Score:      story.Score,
				Blurb:      blurb,
			})
		}
	}

	return renderDigestHTML(digestDate, digestItemsToEmailStories(items), "https://hn30.eu"), nil
}
