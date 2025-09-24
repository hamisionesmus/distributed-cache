package main

import (
	"container/list"
	"sync"
	"time"
)

// CacheEntry represents a cache entry with TTL
type CacheEntry struct {
	Key        string
	Value      []byte
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	AccessCount int64
	LastAccessed time.Time
	element    *list.Element
}

// Cache implements an LRU cache with TTL support
type Cache struct {
	data     map[string]*CacheEntry
	lru      *list.List
	maxSize  int
	currentSize int
	mutex    sync.RWMutex
}

// NewCache creates a new cache with the specified maximum size
func NewCache(maxSize int) *Cache {
	return &Cache{
		data:    make(map[string]*CacheEntry),
		lru:     list.New(),
		maxSize: maxSize,
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		c.removeEntry(entry)
		return nil, false
	}

	// Update access statistics
	entry.AccessCount++
	entry.LastAccessed = time.Now()

	// Move to front (most recently used)
	c.lru.MoveToFront(entry.element)

	return entry.Value, true
}

// Set stores a value in the cache with optional TTL
func (c *Cache) Set(key string, value []byte, ttl *time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Remove existing entry if it exists
	if entry, exists := c.data[key]; exists {
		c.removeEntry(entry)
	}

	// Create new entry
	entry := &CacheEntry{
		Key:         key,
		Value:       value,
		CreatedAt:   time.Now(),
		LastAccessed: time.Now(),
		AccessCount: 0,
	}

	if ttl != nil {
		expiresAt := time.Now().Add(*ttl)
		entry.ExpiresAt = &expiresAt
	}

	// Add to LRU list
	entry.element = c.lru.PushFront(entry)
	c.data[key] = entry
	c.currentSize++

	// Evict if over capacity
	for c.currentSize > c.maxSize && c.lru.Len() > 0 {
		c.evictLRU()
	}
}

// Delete removes a key from the cache
func (c *Cache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, exists := c.data[key]; exists {
		c.removeEntry(entry)
		return true
	}
	return false
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return false
	}

	// Check expiration
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		// Note: We don't remove here to avoid write lock in read operation
		// The entry will be cleaned up on next access
		return false
	}

	return true
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*CacheEntry)
	c.lru = list.New()
	c.currentSize = 0
}

// Stats returns cache statistics
func (c *Cache) Stats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	totalAccesses := int64(0)
	totalSize := 0

	for _, entry := range c.data {
		totalAccesses += entry.AccessCount
		totalSize += len(entry.Value)
	}

	return map[string]interface{}{
		"total_keys":     len(c.data),
		"max_size":       c.maxSize,
		"current_size":   c.currentSize,
		"total_accesses": totalAccesses,
		"total_size_bytes": totalSize,
		"hit_rate":       c.calculateHitRate(),
	}
}

// Cleanup removes expired entries
func (c *Cache) Cleanup() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expired := 0
	for key, entry := range c.data {
		if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
			c.removeEntry(entry)
			delete(c.data, key)
			expired++
		}
	}

	return expired
}

func (c *Cache) removeEntry(entry *CacheEntry) {
	c.lru.Remove(entry.element)
	delete(c.data, entry.Key)
	c.currentSize--
}

func (c *Cache) evictLRU() {
	element := c.lru.Back()
	if element != nil {
		entry := element.Value.(*CacheEntry)
		c.removeEntry(entry)
	}
}

func (c *Cache) calculateHitRate() float64 {
	totalRequests := int64(0)
	totalHits := int64(0)

	for _, entry := range c.data {
		totalRequests += entry.AccessCount
		if entry.AccessCount > 0 {
			totalHits++
		}
	}

	if totalRequests == 0 {
		return 0.0
	}

	return float64(totalHits) / float64(len(c.data))
}

// StartCleanupRoutine starts a background cleanup routine
func (c *Cache) StartCleanupRoutine(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			expired := c.Cleanup()
			if expired > 0 {
				// Could add logging here
			}
		}
	}()
}