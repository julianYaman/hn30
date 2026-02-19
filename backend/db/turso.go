package db

import (
	"database/sql"
	"hn30/backend/types"
	"log/slog"
	"net/url"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var tursoConn *sql.DB

func OpenTurso(dbURL, authToken string) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "turso_operation",
		"operation", "open",
	)
	start := time.Now()

	logger.Info("opening turso connection",
		"event", "turso_open_started",
		"db_url", dbURL,
	)

	connStr := dbURL + "?authToken=" + authToken
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		logger.Error("turso connection failed",
			"event", "turso_open_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		logger.Error("turso ping failed",
			"event", "turso_ping_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	tursoConn = db

	logger.Info("turso connection established",
		"event", "turso_open_completed",
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func CloseTurso() error {
	if tursoConn != nil {
		return tursoConn.Close()
	}
	return nil
}

func SyncTopStories(stories []types.Story) error {
	if tursoConn == nil {
		return nil
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "turso_operation",
		"operation", "sync_top_stories",
		"story_count", len(stories),
	)
	start := time.Now()

	logger.Info("syncing top stories to turso",
		"event", "turso_sync_started",
	)

	tx, err := tursoConn.Begin()
	if err != nil {
		logger.Error("turso transaction start failed",
			"event", "turso_tx_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT hn_id FROM top_stories")
	if err != nil {
		logger.Error("turso query failed",
			"event", "turso_query_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	existingIDs := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			logger.Error("turso scan failed",
				"event", "turso_scan_failed",
				"error", err,
				"duration_ms", time.Since(start).Milliseconds(),
			)
			return err
		}
		existingIDs[id] = true
	}
	rows.Close()

	newIDs := make(map[int]bool)
	for _, s := range stories {
		newIDs[s.ID] = true
	}

	deletedCount := 0
	for id := range existingIDs {
		if !newIDs[id] {
			_, err = tx.Exec("DELETE FROM top_stories WHERE hn_id = ?", id)
			if err != nil {
				logger.Error("turso delete failed",
					"event", "turso_delete_failed",
					"story_id", id,
					"error", err,
					"duration_ms", time.Since(start).Milliseconds(),
				)
				return err
			}
			deletedCount++
		}
	}

	now := time.Now().Unix()
	insertedCount := 0
	updatedCount := 0

	for _, s := range stories {
		domain := extractDomain(s.URL)
		if existingIDs[s.ID] {
			_, err = tx.Exec(
				"UPDATE top_stories SET title = ?, domain = ? WHERE hn_id = ?",
				s.Title, domain, s.ID,
			)
			if err != nil {
				logger.Error("turso update failed",
					"event", "turso_update_failed",
					"story_id", s.ID,
					"error", err,
					"duration_ms", time.Since(start).Milliseconds(),
				)
				return err
			}
			updatedCount++
		} else {
			_, err = tx.Exec(
				"INSERT INTO top_stories (hn_id, title, domain, added_at) VALUES (?, ?, ?, ?)",
				s.ID, s.Title, domain, now,
			)
			if err != nil {
				logger.Error("turso insert failed",
					"event", "turso_insert_failed",
					"story_id", s.ID,
					"error", err,
					"duration_ms", time.Since(start).Milliseconds(),
				)
				return err
			}
			insertedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("turso commit failed",
			"event", "turso_commit_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	logger.Info("turso sync completed",
		"event", "turso_sync_completed",
		"deleted_count", deletedCount,
		"inserted_count", insertedCount,
		"updated_count", updatedCount,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func extractDomain(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	return u.Host
}
