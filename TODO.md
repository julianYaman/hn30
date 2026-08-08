# Newsletter — remaining work

## Done

- SQLite: `story_blurbs`, `digests`, `digest_items` (+ claim / sent / failed)
- Short AI digest blurbs at send time (skip on extract/AI failure; blocked hosts like x.com)
- HTML email template + **Read in browser** → `https://hn30.eu/digest`
- Public website page: `/digest` serves the same HTML as the email (via SvelteKit wrapper; no site chrome)
- SvelteKit API proxy: `/api/newsletter/digest` → backend
- Mailjet campaign send + contact list subscribe
- Scheduler: **07:30 Europe/Berlin** when `NEWSLETTER_ENABLED=true`
- Daily dispatch subscribe form wired
- **Removed** `generate=1`, navbar link, and custom styled digest page

## Env vars — **api-hn-news** (backend) only

```env
MAILJET_API_KEY=
MAILJET_API_SECRET=
MAILJET_LIST_ID=
MAILJET_SENDER_EMAIL=
MAILJET_SENDER_NAME=hn30

NEWSLETTER_ENABLED=false

OPENROUTER_DIGEST_MODEL=@preset/hn30-digest-blurb
OPENROUTER_API_KEY=
SQLITE_PATH=./data/hn30.db
```

Frontend (**hn-news**) needs no Mailjet/newsletter secrets (only existing `PRIVATE_API_BASE_URL`).

`NEWSLETTER_PREVIEW_TOKEN` is no longer used — remove it if you added it.

## Your checklist

1. Mailjet list + verified sender + env vars on **api-hn-news**
2. OpenRouter digest preset (or override model)
3. Persistent SQLite volume (already confirmed)
4. Subscribe yourself only → keep `NEWSLETTER_ENABLED=false` until ready
5. First live send: `NEWSLETTER_ENABLED=true` (emails the whole list)
6. Confirm Mailjet unsub tag `[[UNSUB_LINK_EN]]`

## Follow-ups

- [ ] Manual admin trigger to force-build/send a digest (testing)
- [ ] Pre-generate blurbs before 07:30 to avoid a morning burst
- [ ] Rate-limit subscribe endpoint
- [ ] Optional: date archives `/digest/[date]`
- [ ] Remove local `newsletter-preview.html` from deploys if unused

## API

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/newsletter/digest` | Public email HTML (SvelteKit-proxied; never generates) |
| GET | `/digest` | Same HTML at a public page URL |
| POST | `/api/newsletter/subscribe` | `{ "email": "…" }` → Mailjet list |
