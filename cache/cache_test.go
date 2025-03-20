package cache

import (
	"sync"
	"testing"
	"time"
)

func TestCacheBasicOperations(t *testing.T) {
	cache := NewInMemoryCache()

	// Test Set and Get
	cache.Set("key1", "value1")
	value, found := cache.Get("key1")
	if !found {
		t.Errorf("Expected to find key1 in cache, but it was not found")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// Test Get non-existent key
	_, found = cache.Get("nonexistent")
	if found {
		t.Errorf("Did not expect to find nonexistent key in cache")
	}

	// Test Delete
	cache.Delete("key1")
	_, found = cache.Get("key1")
	if found {
		t.Errorf("Expected key1 to be deleted, but it was found")
	}

	// Test Clear
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Clear()
	_, found = cache.Get("key1")
	if found {
		t.Errorf("Expected cache to be cleared, but key1 was found")
	}
	_, found = cache.Get("key2")
	if found {
		t.Errorf("Expected cache to be cleared, but key2 was found")
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := NewInMemoryCache()

	// Set an item with a short expiration
	cache.SetWithExpiration("key1", "value1", 100*time.Millisecond)

	// Item should exist immediately
	_, found := cache.Get("key1")
	if !found {
		t.Errorf("Expected to find key1 in cache immediately after setting")
	}

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	// Item should be expired now
	_, found = cache.Get("key1")
	if found {
		t.Errorf("Expected key1 to be expired and not found")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewInMemoryCache()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2) // For both readers and writers

	// Start multiple writers
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			cache.Set(string([]rune{'A' + rune(i)}), i)
		}(i)
	}

	// Start multiple readers
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			// Just read value, we don't care about actual value
			cache.Get(string([]rune{'A' + rune(i%goroutines)}))
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// If we get here without panics, the test passes
}
