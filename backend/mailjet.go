package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const mailjetAPIBase = "https://api.mailjet.com/v3/REST"

type mailjetClient struct {
	apiKey      string
	apiSecret   string
	listID      int
	senderID    int
	senderName  string
	senderEmail string
	http        *http.Client
}

func newMailjetClientFromEnv() (*mailjetClient, error) {
	apiKey := strings.TrimSpace(os.Getenv("MAILJET_API_KEY"))
	apiSecret := strings.TrimSpace(os.Getenv("MAILJET_API_SECRET"))
	listIDStr := strings.TrimSpace(os.Getenv("MAILJET_LIST_ID"))
	senderIDStr := strings.TrimSpace(os.Getenv("MAILJET_SENDER_ID"))
	senderName := strings.TrimSpace(os.Getenv("MAILJET_SENDER_NAME"))
	senderEmail := strings.TrimSpace(os.Getenv("MAILJET_SENDER_EMAIL"))

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("MAILJET_API_KEY and MAILJET_API_SECRET are required")
	}
	if listIDStr == "" {
		return nil, fmt.Errorf("MAILJET_LIST_ID is required")
	}
	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAILJET_LIST_ID: %w", err)
	}
	if senderIDStr == "" {
		return nil, fmt.Errorf("MAILJET_SENDER_ID is required")
	}
	senderID, err := strconv.Atoi(senderIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAILJET_SENDER_ID: %w", err)
	}
	if senderEmail == "" {
		return nil, fmt.Errorf("MAILJET_SENDER_EMAIL is required")
	}
	if senderName == "" {
		senderName = "hn30 Daily Dispatch"
	}

	return &mailjetClient{
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		listID:      listID,
		senderID:    senderID,
		senderName:  senderName,
		senderEmail: senderEmail,
		http:        &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *mailjetClient) doJSON(method, path string, payload any, dest any) error {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, mailjetAPIBase+path, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("mailjet %s %s: status %d: %s", method, path, res.StatusCode, truncate(string(resBody), 500))
	}
	if dest == nil {
		return nil
	}
	if len(resBody) == 0 {
		return nil
	}
	return json.Unmarshal(resBody, dest)
}

type mailjetDataResponse struct {
	Count int                      `json:"Count"`
	Data  []map[string]interface{} `json:"Data"`
}

func (c *mailjetClient) CreateAndSendCampaign(title, subject, htmlPart, textPart string) (int64, error) {
	var createResp mailjetDataResponse
	err := c.doJSON("POST", "/campaigndraft", map[string]any{
		"Locale":         "en_US",
		"Sender":         strconv.Itoa(c.senderID),
		"SenderEmail":    c.senderEmail,
		"SenderName":     c.senderName,
		"Subject":        subject,
		"ContactsListID": c.listID,
		"Title":          title,
	}, &createResp)
	if err != nil {
		return 0, fmt.Errorf("create campaign draft: %w", err)
	}
	if len(createResp.Data) == 0 {
		return 0, fmt.Errorf("create campaign draft: empty response")
	}

	draftID, ok := asInt64(createResp.Data[0]["ID"])
	if !ok {
		return 0, fmt.Errorf("create campaign draft: missing ID")
	}

	err = c.doJSON("POST", fmt.Sprintf("/campaigndraft/%d/detailcontent", draftID), map[string]any{
		"Html-part": htmlPart,
		"Text-part": textPart,
	}, nil)
	if err != nil {
		return draftID, fmt.Errorf("set campaign content: %w", err)
	}

	err = c.doJSON("POST", fmt.Sprintf("/campaigndraft/%d/send", draftID), nil, nil)
	if err != nil {
		return draftID, fmt.Errorf("send campaign: %w", err)
	}

	return draftID, nil
}

func (c *mailjetClient) SubscribeEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email")
	}

	path := fmt.Sprintf("/contactslist/%d/managecontact", c.listID)
	return c.doJSON("POST", path, map[string]any{
		"Email":  email,
		"Action": "addnoforce",
	}, nil)
}

func subscribeErrorStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	msg := err.Error()
	if strings.Contains(msg, "invalid email") {
		return http.StatusBadRequest
	}
	return http.StatusBadGateway
}

func subscribeErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if strings.Contains(msg, "invalid email") {
		return "Invalid email address"
	}
	if strings.Contains(msg, "Object not found") || strings.Contains(msg, "status 404") {
		return "Newsletter list not found. Check MAILJET_LIST_ID matches a list in the same Mailjet account as your API keys."
	}
	if strings.Contains(msg, "MAILJET_") || strings.Contains(msg, "required") {
		return "Newsletter subscribe is not configured"
	}
	return "Could not subscribe email"
}

func asInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int64:
		return n, true
	case int:
		return int64(n), true
	case json.Number:
		i, err := n.Int64()
		return i, err == nil
	default:
		return 0, false
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
