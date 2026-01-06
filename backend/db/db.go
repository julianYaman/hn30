package db

import (
	"database/sql"
	"fmt"
	"hn30/backend/types"
	"hn30/backend/utils"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) *sql.DB {
	db, err := sql.Open("sqlite", path+"?_busy_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}

	// SQLite pragmas
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Fatal(err)
		}
	}

	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Hour)

	utils.LogComponent("DB", "Database opened at %s", path)

	return db
}

func migrate(db *sql.DB) error {
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
	return err
}

func UpsertStory(db *sql.DB, s types.Story) error {
	now := time.Now().Unix()

	_, err := db.Exec(`
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

	utils.LogComponent("DB", "Upserted story %d into database", s.ID)

	return err
}

func ShouldNotify(db *sql.DB, s types.Story) bool {
	var notifiedAt sql.NullInt64
	var createdTime int64
	var maxPoints int

	utils.LogComponent("DB", "Checking story %d", s.ID)

	err := db.QueryRow(`
		SELECT notified_at, created_at, max_points
		FROM stories
		WHERE hn_id = ?
	`, s.ID).Scan(&notifiedAt, &createdTime, &maxPoints)

	fmt.Println("notifiedAt:", notifiedAt, "createdTime:", createdTime, "now:", time.Now().Unix(), "maxPoints:", maxPoints)

	if err != nil {
		return false
	}

	if notifiedAt.Valid {
		fmt.Println("Story has already been notified at:", notifiedAt.Int64)
		return false
	}

	age := time.Now().Unix() - createdTime
	if age < 60*60 {
		fmt.Println("Story is too new, age (s):", age)
		return false
	}

	if maxPoints < 600 {
		fmt.Println("Story does not have enough points:", maxPoints)
		return false
	}

	utils.LogComponent("NOTIFICATION", "Story %d is eligible for notification", s.ID)

	return true
}

func MarkNotified(db *sql.DB, storyID int) error {
	_, err := db.Exec(`
		UPDATE stories
		SET notified_at = ?
		WHERE hn_id = ?
	`, time.Now().Unix(), storyID)

	if err == nil {
		utils.LogComponent("DB", "Marked story %d as notified", storyID)
	}

	return err
}
