package main

import (
	"hn30/backend/utils"
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
	for id := range c.stories {
		if !newIDs[id] {
			utils.LogComponent("CACHE", "Removing story ID %d from cache", id)
			delete(c.stories, id)
		}
	}

	c.storyIDs = ids
}

func (c *Cache) GetAll() []EnrichedStory {
	c.mu.RLock()
	defer c.mu.RUnlock()
	stories := make([]EnrichedStory, 0, len(c.storyIDs))
	for _, id := range c.storyIDs {
		story, found := c.stories[id]
		if !found {
			utils.LogComponent("CACHE", "Story ID %d not found in cache", id)
			continue
		}
		stories = append(stories, story)
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
