package gocache_ttl

import (
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	cache := NewCache(time.Second, 10)
	defer cache.Close()

	cache.Set("key", "value", 2*time.Second)

	it, ok := cache.items["key"]
	if !ok {
		t.Fatal("expected key  in cache.items")
	}
	if it.value.(string) != "value" {
		t.Fatalf("expected 'value', got %v", it.value)
	}
	if it.expiration == 0 {
		t.Fatal("expected expiration to be set")
	}
}

func TestCache_Get(t *testing.T) {
	cache := NewCache(time.Second, 10)
	defer cache.Close()

	cache.Set("key", "value", 2*time.Second)

	it, ok := cache.Get("key")
	if !ok {
		t.Fatal("expected key in cache.items")
	}
	if it.(string) != "value" {
		t.Fatalf("expected 'value', got %v", it)
	}
	time.Sleep(2 * time.Second)
	if _, ok := cache.Get("key"); ok {
		t.Fatal("expected key to expire")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(time.Second, 10)
	defer cache.Close()

	cache.Set("key", "value", 2*time.Second)
	cache.Delete("key")
	_, ok := cache.Get("key")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCache_Exists(t *testing.T) {
	cache := NewCache(time.Second, 10)
	defer cache.Close()

	cache.Set("key", "value", 2*time.Second)
	if !cache.Exists("key") {
		t.Fatal("expected key to exist")
	}
	time.Sleep(3 * time.Second)
	if cache.Exists("key") {
		t.Fatal("expected key to not exist")
	}
}

func TestCache_Keys(t *testing.T) {
	cache := NewCache(time.Second, 10)
	defer cache.Close()

	cache.Set("key", "value", 2*time.Second)
	cache.Set("key2", "value2", 2*time.Second)
	cache.Set("key3", "value3", 3*time.Second)

	keys := cache.Keys()

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}

	expectedKeys := map[string]bool{
		"key":  true,
		"key2": true,
		"key3": true,
	}
	for _, key := range keys {
		if !expectedKeys[key] {
			t.Fatalf("unexpected key: %s", key)
		}
	}
}

func TestCache_cleanup(t *testing.T) {
	cache := NewCache(time.Second, 10)
	defer cache.Close()

	cache.Set("key", "value", 2*time.Second)

	if _, ok := cache.Get("key"); !ok {
		t.Fatal("expected key to exist right after Set")
	}

	time.Sleep(3 * time.Second)

	if _, ok := cache.Get("key"); ok {
		t.Fatal("expected key to expire and be cleaned up")
	}
}
