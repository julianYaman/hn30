package main

import (
	"encoding/json"
	"hn30/backend/utils"
	"net/http"
	"strconv"
)

func topStoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Content-Type", "application/json")
	stories := storyCache.GetAll()
	json.NewEncoder(w).Encode(stories)
}

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Content-Type", "application/json")

	// 1. Get and validate the story ID from the query params
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing story ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid story ID", http.StatusBadRequest)
		return
	}

	// 2. Check if the story exists in the cache
	story, found := storyCache.Get(id)
	if !found {
		http.Error(w, "Story not found", http.StatusNotFound)
		return
	}

	// 3. If summary already exists, return it immediately
	if story.Summary != "" {
		utils.LogComponent("CACHE", "Returning cached summary for story %d", id)
		json.NewEncoder(w).Encode(map[string]string{"summary": story.Summary, "model": story.SummaryModel})
		return
	}

	// 4. If no summary, generate one
	utils.LogComponent("SUMMARIZER", "No summary found for story %d, generating...", id)
	articleText, err := extractArticleText(story.URL)
	if err != nil {
		utils.LogError("Failed to extract article text for story %d: %v", id, err)
		http.Error(w, "Failed to extract article content", http.StatusInternalServerError)
		return
	}

	summary, err := generateSummary(articleText)
	if err != nil {
		utils.LogError("Failed to generate summary for story %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Save the new summary and article text to the cache
	story.Summary = summary.Summary
	story.ArticleText = articleText
	story.SummaryModel = summary.Model
	storyCache.Set(id, story)
	utils.LogComponent("CACHE", "Saved new summary for story %d to cache", id)

	// 6. Return the new summary
	json.NewEncoder(w).Encode(map[string]string{"summary": summary.Summary, "model": summary.Model})
}
