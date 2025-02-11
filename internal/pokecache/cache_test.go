package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestCacheConcurrency(t *testing.T) {
    const interval = 5 * time.Second
    cache := NewCache(interval)
    
    // Test concurrent writes
    for i := 0; i < 100; i++ {
        go func(i int) {
            key := fmt.Sprintf("key%d", i)
            cache.Add(key, []byte(fmt.Sprintf("value%d", i)))
        }(i)
    }
    
    // Let all goroutines finish
    time.Sleep(100 * time.Millisecond)
    
    // Verify we can read all values
    for i := 0; i < 100; i++ {
        key := fmt.Sprintf("key%d", i)
        _, ok := cache.Get(key)
        if !ok {
            t.Errorf("failed to get key %s", key)
        }
    }
}

func TestEmptyCache(t *testing.T) {
    cache := NewCache(5 * time.Second)
    _, ok := cache.Get("nonexistent")
    if ok {
        t.Error("expected false for nonexistent key")
    }
}

func TestOverwriteValue(t *testing.T) {
    cache := NewCache(5 * time.Second)
    key := "test-key"
    
    // Add initial value
    cache.Add(key, []byte("initial"))
    
    // Overwrite with new value
    cache.Add(key, []byte("updated"))
    
    val, ok := cache.Get(key)
    if !ok {
        t.Error("expected to find key")
        return
    }
    if string(val) != "updated" {
        t.Errorf("expected 'updated', got '%s'", string(val))
    }
}

func TestMultipleReaps(t *testing.T) {
    // Use longer intervals to make test more reliable
    cache := NewCache(100 * time.Millisecond)
    
    // Add first item
    cache.Add("key1", []byte("value1"))
    
    // Wait a bit and add second item
    time.Sleep(50 * time.Millisecond)
    cache.Add("key2", []byte("value2"))
    
    // Both should exist
    _, ok1 := cache.Get("key1")
    _, ok2 := cache.Get("key2")
    if !ok1 || !ok2 {
        t.Error("expected both keys to exist initially")
        return
    }
    
    // Wait longer than the cache interval to ensure first reap
    time.Sleep(150 * time.Millisecond)
    
    // key1 should be gone, key2 should exist
    _, ok1 = cache.Get("key1")
    _, ok2 = cache.Get("key2")
    if ok1 {
        t.Error("expected key1 to be reaped")
    }
    if !ok2 {
        t.Error("expected key2 to still exist")
    }
    
    // Wait for second reap
    time.Sleep(150 * time.Millisecond)
    
    // Both should be gone
    _, ok1 = cache.Get("key1")
    _, ok2 = cache.Get("key2")
    if ok1 || ok2 {
        t.Error("expected both keys to be reaped")
    }
}

func TestReapEmptyCache(t *testing.T) {
    cache := NewCache(10 * time.Millisecond)
    
    // Add and immediately verify a key doesn't exist
    _, ok := cache.Get("nonexistent")
    if ok {
        t.Error("expected empty cache to return false for any key")
    }
    
    // Wait for multiple reap cycles on empty cache
    time.Sleep(50 * time.Millisecond)
    
    // Verify we can still add and get items after reaping empty cache
    cache.Add("newkey", []byte("value"))
    val, ok := cache.Get("newkey")
    if !ok {
        t.Error("expected to find new key after reap cycles")
    }
    if string(val) != "value" {
        t.Error("expected to get correct value after reap cycles")
    }
}