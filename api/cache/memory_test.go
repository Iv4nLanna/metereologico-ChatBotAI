package cache

import (
	"testing"
	"time"
)

func TestSet_Get(t *testing.T) {
	c := New(1 * time.Minute)
	c.Set("key", "value")
	got, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got != "value" {
		t.Fatalf("expected 'value', got %v", got)
	}
}

func TestGet_Missing(t *testing.T) {
	c := New(1 * time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss for missing key")
	}
}

func TestGet_Expired(t *testing.T) {
	c := New(10 * time.Millisecond)
	c.Set("key", "value")
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected cache miss after TTL expiration")
	}
}
