package db

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	BlurbStatusOK      = "ok"
	BlurbStatusSkipped = "skipped"
	BlurbStatusPending = "pending"

	DigestStatusPending = "pending"
	DigestStatusSent    = "sent"
	DigestStatusFailed  = "failed"
)

type StoryBlurb struct {
	HNID      int
	Blurb     string
	Status    string
	Model     string
	CreatedAt int64
	UpdatedAt int64
}

type Digest struct {
	DigestDate     string
	Status         string
	ClaimedAt      int64
	SentAt         sql.NullInt64
	Error          string
	MailjetDraftID sql.NullInt64
}

type DigestItem struct {
	DigestDate string
	Rank       int
	HNID       int
	Title      string
	URL        string
	Score      int
	Blurb      string
}

// ClaimDigest tries to claim today's digest job.
// Returns true if this process owns the claim and should run the send.
func ClaimDigest(db *sql.DB, digestDate string) (bool, error) {
	now := time.Now().Unix()

	res, err := db.Exec(`
		INSERT INTO digests (digest_date, status, claimed_at, error)
		VALUES (?, ?, ?, '')
		ON CONFLICT(digest_date) DO UPDATE SET
			status = excluded.status,
			claimed_at = excluded.claimed_at,
			error = ''
		WHERE digests.status != ?
	`, digestDate, DigestStatusPending, now, DigestStatusSent)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func MarkDigestSent(db *sql.DB, digestDate string, mailjetDraftID int64) error {
	now := time.Now().Unix()
	_, err := db.Exec(`
		UPDATE digests
		SET status = ?, sent_at = ?, mailjet_draft_id = ?, error = ''
		WHERE digest_date = ?
	`, DigestStatusSent, now, mailjetDraftID, digestDate)
	return err
}

func MarkDigestFailed(db *sql.DB, digestDate string, digesterr error) error {
	msg := ""
	if digesterr != nil {
		msg = digesterr.Error()
		if len(msg) > 1000 {
			msg = msg[:1000]
		}
	}
	_, err := db.Exec(`
		UPDATE digests
		SET status = ?, error = ?
		WHERE digest_date = ?
	`, DigestStatusFailed, msg, digestDate)
	return err
}

func GetDigest(db *sql.DB, digestDate string) (*Digest, error) {
	var d Digest
	err := db.QueryRow(`
		SELECT digest_date, status, claimed_at, sent_at, error, mailjet_draft_id
		FROM digests
		WHERE digest_date = ?
	`, digestDate).Scan(
		&d.DigestDate, &d.Status, &d.ClaimedAt, &d.SentAt, &d.Error, &d.MailjetDraftID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func ReplaceDigestItems(db *sql.DB, digestDate string, items []DigestItem) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM digest_items WHERE digest_date = ?`, digestDate); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO digest_items (digest_date, rank, hn_id, title, url, score, blurb)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err := stmt.Exec(
			digestDate, item.Rank, item.HNID, item.Title, item.URL, item.Score, item.Blurb,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetDigestItems(db *sql.DB, digestDate string) ([]DigestItem, error) {
	rows, err := db.Query(`
		SELECT digest_date, rank, hn_id, title, url, score, blurb
		FROM digest_items
		WHERE digest_date = ?
		ORDER BY rank ASC
	`, digestDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]DigestItem, 0)
	for rows.Next() {
		var item DigestItem
		if err := rows.Scan(
			&item.DigestDate, &item.Rank, &item.HNID, &item.Title, &item.URL, &item.Score, &item.Blurb,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func GetStoryBlurb(db *sql.DB, hnID int) (*StoryBlurb, error) {
	var b StoryBlurb
	err := db.QueryRow(`
		SELECT hn_id, blurb, status, model, created_at, updated_at
		FROM story_blurbs
		WHERE hn_id = ?
	`, hnID).Scan(&b.HNID, &b.Blurb, &b.Status, &b.Model, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func UpsertStoryBlurb(db *sql.DB, hnID int, blurb, status, model string) error {
	now := time.Now().Unix()
	_, err := db.Exec(`
		INSERT INTO story_blurbs (hn_id, blurb, status, model, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(hn_id) DO UPDATE SET
			blurb = excluded.blurb,
			status = excluded.status,
			model = excluded.model,
			updated_at = excluded.updated_at
	`, hnID, blurb, status, model, now, now)
	return err
}

func BerlinDigestDate(t time.Time) (string, error) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return "", fmt.Errorf("load Europe/Berlin: %w", err)
	}
	return t.In(loc).Format("2006-01-02"), nil
}
