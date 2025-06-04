package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	key := "testKey"
	val := []byte("testValue")

	cache.Add(key, val)

	retrievedVal, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key %s, but it was not found", key)
	}
	if string(retrievedVal) != string(val) {
		t.Errorf("expected value %s, but got %s", string(val), string(retrievedVal))
	}
}

func TestCacheGetNotFound(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	key := "nonExistentKey"

	_, ok := cache.Get(key)
	if ok {
		t.Errorf("expected not to find key %s, but it was found", key)
	}
}

func TestCacheReap(t *testing.T) {
	const reapInterval = 50 * time.Millisecond // Short interval for testing
	const testDelay = 100 * time.Millisecond   // Must be > reapInterval

	cache := NewCache(reapInterval)

	key1 := "key1"
	val1 := []byte("val1")
	cache.Add(key1, val1)

	// Wait for a period longer than the reap interval
	time.Sleep(testDelay)

	_, ok := cache.Get(key1)
	if ok {
		t.Errorf("expected key %s to be reaped, but it was found", key1)
	}
}

func TestCacheReapLoop(t *testing.T) {
	const reapInterval = 50 * time.Millisecond
	const testDelay = 30 * time.Millisecond // Less than reap interval to check if it's NOT reaped too early

	cache := NewCache(reapInterval)

	key := "testKeyActive"
	val := []byte("testValueActive")
	cache.Add(key, val)

	// Wait for a period shorter than the reap interval
	time.Sleep(testDelay)

	// Item should still be in cache
	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected key %s to still be in cache before reap, but it was not found", key)
	}

	// Wait for longer than the reap interval to ensure it's reaped
	time.Sleep(reapInterval * 2)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected key %s to be reaped after interval, but it was found", key)
	}
}
