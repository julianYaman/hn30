package main

import (
	"bytes"
	"fmt"
	"html"
	"hn30/backend/db"
	"net/url"
	"strings"
	"time"
)

type digestEmailStory struct {
	Rank   int
	Title  string
	URL    string
	Source string
	Blurb  string
}

func sourceLabel(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	if host == "" {
		return "news.ycombinator.com"
	}

	if host == "github.com" {
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		clean := make([]string, 0, 2)
		for _, p := range parts {
			if p != "" {
				clean = append(clean, p)
			}
			if len(clean) == 2 {
				break
			}
		}
		switch len(clean) {
		case 0:
			return "github.com"
		case 1:
			return "github.com/" + clean[0]
		default:
			label := "github.com/" + clean[0] + "/" + clean[1]
			if len(label) > 48 {
				return label[:47] + "…"
			}
			return label
		}
	}

	return host
}

func digestItemsToEmailStories(items []db.DigestItem) []digestEmailStory {
	out := make([]digestEmailStory, 0, len(items))
	for _, item := range items {
		out = append(out, digestEmailStory{
			Rank:   item.Rank,
			Title:  item.Title,
			URL:    item.URL,
			Source: sourceLabel(item.URL),
			Blurb:  strings.TrimSpace(item.Blurb),
		})
	}
	return out
}

func formatDigestDateLabel(digestDate string) string {
	t, err := time.Parse("2006-01-02", digestDate)
	if err != nil {
		return digestDate
	}
	return t.Format("2 Jan 2006")
}

func renderDigestHTML(digestDate string, stories []digestEmailStory, unsubscribeURL string) string {
	if unsubscribeURL == "" {
		// Mailjet campaign merge tag
		unsubscribeURL = "[[UNSUB_LINK_EN]]"
	}

	var b bytes.Buffer
	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>hn30 daily dispatch</title>
</head>
<body style="margin:0;padding:0;background:#ffffff;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background:#ffffff;">
<tr><td align="center" style="background:#ffffff;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:600px;margin:0 auto;background:#ffffff;">
  <tr>
    <td align="center" style="padding:28px 24px 18px 24px;border-bottom:2px solid #111111;text-align:center;background:#ffffff;">
      <a href="https://hn30.eu" style="font-family:Georgia,'Times New Roman',serif;font-size:28px;line-height:1.1;color:#111111;font-weight:700;text-decoration:none;">hn<span style="color:#ff6600;">30</span></a>
      <div style="margin-top:6px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;font-size:12px;letter-spacing:0.06em;text-transform:uppercase;color:#999999;text-align:center;">Daily dispatch · `)
	b.WriteString(html.EscapeString(formatDigestDateLabel(digestDate)))
	b.WriteString(`</div>
      <div style="margin-top:10px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;font-size:13px;line-height:1.4;text-align:center;">
        <a href="https://hn30.eu/digest" style="color:#ff6600;text-decoration:underline;">Read in browser</a>
      </div>
    </td>
  </tr>
  <tr>
    <td style="padding:8px 24px 8px 24px;background:#ffffff;">
      <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
`)

	for _, story := range stories {
		b.WriteString(`
<tr>
  <td style="padding:14px 0;border-bottom:1px solid #e5e5e5;">
    <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
      <tr>
        <td width="28" valign="top" style="font-family:Georgia,'Times New Roman',serif;font-size:15px;line-height:1.35;color:#999999;padding-right:8px;">`)
		b.WriteString(fmt.Sprintf("%d.", story.Rank))
		b.WriteString(`</td>
        <td valign="top">
          <a href="`)
		b.WriteString(html.EscapeString(story.URL))
		b.WriteString(`" style="font-family:Georgia,'Times New Roman',serif;font-size:16px;line-height:1.35;color:#111111;font-weight:700;text-decoration:none;">`)
		b.WriteString(html.EscapeString(story.Title))
		b.WriteString(`</a>`)
		if story.Source != "" {
			b.WriteString(` <span style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;font-size:12px;line-height:1.35;color:#999999;">(`)
			b.WriteString(html.EscapeString(story.Source))
			b.WriteString(`)</span>`)
		}
		if story.Blurb != "" {
			b.WriteString(`
          <div style="margin:4px 0 0 0;font-family:Georgia,'Times New Roman',serif;font-size:14px;line-height:1.45;color:#555555;">`)
			b.WriteString(html.EscapeString(story.Blurb))
			b.WriteString(`</div>`)
		}
		b.WriteString(`
        </td>
      </tr>
    </table>
  </td>
</tr>
`)
	}

	b.WriteString(`
      </table>
    </td>
  </tr>
  <tr>
    <td style="padding:24px;border-top:1px solid #e5e5e5;background:#ffffff;">
      <div style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;font-size:13px;line-height:1.5;color:#555555;">
        <a href="https://hn30.eu" style="color:#ff6600;text-decoration:none;">Read on hn30.eu</a>
        &nbsp;·&nbsp;
        <a href="`)
	b.WriteString(html.EscapeString(unsubscribeURL))
	b.WriteString(`" style="color:#999999;text-decoration:underline;">Unsubscribe</a>
      </div>
    </td>
  </tr>
</table>
</td></tr>
</table>
</body>
</html>
`)
	return b.String()
}

func renderDigestText(digestDate string, stories []digestEmailStory) string {
	var b strings.Builder
	b.WriteString("hn30 — Daily dispatch · ")
	b.WriteString(formatDigestDateLabel(digestDate))
	b.WriteString("\n\n")
	for _, story := range stories {
		b.WriteString(fmt.Sprintf("%d. %s", story.Rank, story.Title))
		if story.Source != "" {
			b.WriteString(" (")
			b.WriteString(story.Source)
			b.WriteString(")")
		}
		b.WriteString("\n")
		b.WriteString(story.URL)
		b.WriteString("\n")
		if story.Blurb != "" {
			b.WriteString(story.Blurb)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("Read on https://hn30.eu\n")
	return b.String()
}
