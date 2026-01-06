package main

import (
	"log/slog"
	"os"
	"sync"
	"time"
)

type Cache struct {
	mu          sync.RWMutex
	stories     map[int]EnrichedStory
	storyIDs    []int
	lastUpdated time.Time
}

func NewCache() *Cache {
	return &Cache{
		stories:  make(map[int]EnrichedStory),
		storyIDs: make([]int, 0),
	}
}

func (c *Cache) Set(id int, story EnrichedStory) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stories[id] = story
}

func (c *Cache) Get(id int) (EnrichedStory, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	story, found := c.stories[id]
	return story, found
}

func (c *Cache) SetStoryIDs(ids []int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a map for quick lookup of new story IDs
	newIDs := make(map[int]bool)
	for _, id := range ids {
		newIDs[id] = true
	}

	// Remove stories that are no longer in the top 30
	removedCount := 0
	removedIDs := make([]int, 0)
	for id := range c.stories {
		if !newIDs[id] {
			removedIDs = append(removedIDs, id)
			delete(c.stories, id)
			removedCount++
		}
	}

	if removedCount > 0 {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
			"event_type", "cache_operation",
			"operation", "cleanup",
		)
		logger.Info("cache cleanup completed",
			"event", "cache_cleanup",
			"removed_count", removedCount,
			"removed_story_ids", removedIDs,
			"new_story_count", len(ids),
			"remaining_cached_count", len(c.stories),
		)
	}

	c.storyIDs = ids
}

func (c *Cache) GetAll() []EnrichedStory {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stories := make([]EnrichedStory, 0, len(c.storyIDs))
	missingIDs := make([]int, 0)

	for _, id := range c.storyIDs {
		story, found := c.stories[id]
		if !found {
			missingIDs = append(missingIDs, id)
			continue
		}
		stories = append(stories, story)
	}

	// Only log if there are missing stories (cache inconsistency)
	if len(missingIDs) > 0 {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
			"event_type", "cache_operation",
			"operation", "get_all",
		)
		logger.Warn("cache miss for requested stories",
			"event", "cache_miss",
			"missing_story_ids", missingIDs,
			"missing_count", len(missingIDs),
			"returned_count", len(stories),
			"requested_count", len(c.storyIDs),
		)
	}

	return stories
}

func (c *Cache) SetLastUpdated(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastUpdated = t
}

func (c *Cache) LastUpdated() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUpdated
}
