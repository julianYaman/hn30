<p align="center">
  <img src="./frontend/static/favicon.svg" alt="HN News Page Logo" width="120">
</p>

<h1 align="center">hn30 - a yamanlabs project</h1>

hn30 shows the top 30 posts on Hacker News on a tech news-like website. It fetches top stories from the Hacker News API, enriches them with Open Graph metadata, and displays them in a clean, responsive interface with AI-powered features.

The motivation behind this project was to fully utilize vibe-code tools and experiment with them. Also, I wanted to use Gemini for a practical application that I would use daily and might be useful for others as well. Don't expect to see any clean code or best practices here - this was a fun experiment and a learning experience.

## ✨ Features

- **Top 30 Stories:** Displays the current top 30 articles from Hacker News.
- **AI-Powered Summaries:** On-demand summaries of articles generated with OpenRouter.
- **Dynamic Story Data:** Story scores and comment counts are updated in the background.
- **Metadata Enrichment:** Scrapes `og:image` and `og:description` from each article for a richer preview.
- **High-Performance Backend:** A Go microservice with an in-memory cache serves all data instantly.
- **Background Refresh:** The cache is automatically updated every 15 minutes.
- **Bookmarks with Notifications:** Save your favorite stories locally in your browser.

## 🛠️ Tech Stack

- **Backend:** Go
- **Frontend:** SvelteKit, TailwindCSS
- **Reverse Proxy:** Traefik
- **AI Integration:** OpenRouter API

## 🚀 Running the Project

### Running Locally for Development

This method is best for active development, as it provides hot-reloading for the frontend.
A script is provided to run both the backend and frontend concurrently.

**Prerequisites:**
- Go (1.23+)
- Node.js (18+)

**Setup:**

1.  **Environment:**
    The backend requires an OpenRouter API key to generate summaries and the frontend requires the base API url. Create a file named `.env.development` in the `/` directory:
    ```env
    # /.env.development
    OPENROUTER_API_KEY=your_openrouter_api_key_here
    PRIVATE_API_BASE_URL=http://localhost:8080
    ```

2.  **Adjust the AI model you want to use:**
    In [`backend/summarizer.go`](https://github.com/julianYaman/hn30/blob/main/backend/summarizer.go), adjust the model used in the `OpenRouterRequest` with your model that you would like to use.
    ```go
    body := OpenRouterRequest{
  		Model: "<your model>",
  		Messages: []OpenRouterMessage{
  			{
  				Role:    "user",
  				Content: articleText,
  			},
  		},
	  }
    ```   

    For testing purposes, we recommend to use the `:free` models on OpenRouter. In production, we recommend using presets in which you can adjust the used models later on.
    
2.  **Run the development server:**
    Just run the following command in the root directory:
    ```bash
    ./dev.sh
    ```
