package db

import (
	"database/sql"
	"hn30/backend/types"
	"log"
	"log/slog"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) *sql.DB {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "database_operation",
		"operation", "open",
		"db_path", path,
	)
	start := time.Now()

	logger.Info("opening database",
		"event", "db_open_started",
		"connection_string", path+"?_busy_timeout=5000",
	)

	db, err := sql.Open("sqlite", path+"?_busy_timeout=5000")
	if err != nil {
		logger.Error("database open failed",
			"event", "db_open_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		log.Fatal(err)
	}

	logger.Info("database connection opened",
		"event", "db_connection_opened",
		"duration_ms", time.Since(start).Milliseconds(),
	)

	// SQLite pragmas
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
	}

	logger.Info("applying sqlite pragmas",
		"event", "pragmas_apply_started",
		"pragma_count", len(pragmas),
	)

	for i, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			logger.Error("pragma execution failed",
				"event", "pragma_failed",
				"pragma", p,
				"pragma_index", i,
				"error", err,
			)
			log.Fatal(err)
		}
	}

	logger.Info("pragmas applied successfully",
		"event", "pragmas_applied",
		"pragma_count", len(pragmas),
	)

	migrateStart := time.Now()
	logger.Info("running database migrations",
		"event", "migration_started",
	)

	if err := migrate(db); err != nil {
		logger.Error("migration failed",
			"event", "migration_failed",
			"error", err,
			"duration_ms", time.Since(migrateStart).Milliseconds(),
		)
		log.Fatal(err)
	}

	logger.Info("migrations completed",
		"event", "migration_completed",
		"duration_ms", time.Since(migrateStart).Milliseconds(),
	)

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Hour)

	logger.Info("database ready",
		"event", "db_open_completed",
		"max_open_conns", 1,
		"conn_max_lifetime", "1h",
		"total_duration_ms", time.Since(start).Milliseconds(),
	)

	return db
}

func migrate(db *sql.DB) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "database_migration",
	)
	start := time.Now()

	logger.Info("executing migration schema",
		"event", "schema_execution_started",
	)

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS stories (
			hn_id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			url TEXT,
			created_at INTEGER NOT NULL,
			last_seen_at INTEGER NOT NULL,
			max_points INTEGER NOT NULL,
			notified_at INTEGER
		);

		CREATE INDEX IF NOT EXISTS idx_notified
		ON stories (notified_at);
	`)

	if err != nil {
		logger.Error("schema execution failed",
			"event", "schema_execution_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	logger.Info("migration schema applied",
		"event", "schema_execution_completed",
		"tables", []string{"stories"},
		"indexes", []string{"idx_notified"},
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func UpsertStory(db *sql.DB, s types.Story) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "database_operation",
		"operation", "upsert_story",
		"story_id", s.ID,
		"story_title", s.Title,
	)
	start := time.Now()
	now := time.Now().Unix()

	logger.Info("upserting story",
		"event", "upsert_started",
		"story_url", s.URL,
		"story_score", s.Score,
		"story_time", s.Time,
	)

	result, err := db.Exec(`
		INSERT INTO stories (
			hn_id, title, url,
			created_at, last_seen_at,
			max_points
		) VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(hn_id) DO UPDATE SET
			title = excluded.title,
			url = excluded.url,
			last_seen_at = ?,
			max_points = MAX(max_points, excluded.max_points)
		`,
		s.ID, s.Title, s.URL,
		s.Time, now,
		s.Score,
		now,
	)

	if err != nil {
		logger.Error("upsert failed",
			"event", "upsert_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId()

	logger.Info("story upserted successfully",
		"event", "upsert_completed",
		"rows_affected", rowsAffected,
		"last_insert_id", lastInsertId,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func ShouldNotify(db *sql.DB, s types.Story) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "notification_check",
		"story_id", s.ID,
		"story_title", s.Title,
		"story_score", s.Score,
	)
	start := time.Now()

	var notifiedAt sql.NullInt64
	var createdTime int64
	var maxPoints int

	logger.Info("checking notification eligibility",
		"event", "notification_check_started",
	)

	err := db.QueryRow(`
		SELECT notified_at, created_at, max_points
		FROM stories
		WHERE hn_id = ?
	`, s.ID).Scan(&notifiedAt, &createdTime, &maxPoints)

	if err != nil {
		logger.Warn("notification check query failed",
			"event", "query_failed",
			"error", err,
			"eligible", false,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return false
	}

	now := time.Now().Unix()
	age := now - createdTime

	logger.Info("story data retrieved",
		"event", "data_retrieved",
		"notified_at", notifiedAt.Int64,
		"notified_at_valid", notifiedAt.Valid,
		"created_at", createdTime,
		"max_points", maxPoints,
		"age_seconds", age,
		"current_time", now,
	)

	if notifiedAt.Valid {
		logger.Info("notification already sent",
			"event", "notification_check_completed",
			"reason", "already_notified",
			"notified_at_timestamp", notifiedAt.Int64,
			"eligible", false,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return false
	}

	if age < 60*60 {
		logger.Info("story too new",
			"event", "notification_check_completed",
			"reason", "too_new",
			"age_seconds", age,
			"required_age_seconds", 3600,
			"eligible", false,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return false
	}

	if maxPoints < 600 {
		logger.Info("insufficient points",
			"event", "notification_check_completed",
			"reason", "insufficient_points",
			"max_points", maxPoints,
			"required_points", 600,
			"eligible", false,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return false
	}

	logger.Info("story eligible for notification",
		"event", "notification_check_completed",
		"reason", "eligible",
		"age_seconds", age,
		"max_points", maxPoints,
		"eligible", true,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return true
}

func MarkNotified(db *sql.DB, storyID int) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"event_type", "database_operation",
		"operation", "mark_notified",
		"story_id", storyID,
	)
	start := time.Now()
	notifiedAt := time.Now().Unix()

	logger.Info("marking story as notified",
		"event", "mark_notified_started",
		"notified_at_timestamp", notifiedAt,
	)

	result, err := db.Exec(`
		UPDATE stories
		SET notified_at = ?
		WHERE hn_id = ?
	`, notifiedAt, storyID)

	if err != nil {
		logger.Error("mark notified failed",
			"event", "mark_notified_failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	logger.Info("story marked as notified",
		"event", "mark_notified_completed",
		"rows_affected", rowsAffected,
		"notified_at_timestamp", notifiedAt,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}
