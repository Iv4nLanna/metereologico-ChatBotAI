package cache

import (
	"sync"
	"time"

	"weather-api/metrics"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Memory struct {
	mu    sync.RWMutex
	store map[string]entry
	ttl   time.Duration
}

func New(ttl time.Duration) *Memory {
	m := &Memory{
		store: make(map[string]entry),
		ttl:   ttl,
	}
	go m.gc()
	return m
}

func (m *Memory) Set(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = entry{value: value, expiresAt: time.Now().Add(m.ttl)}
}

func (m *Memory) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.store[key]
	if !ok || time.Now().After(e.expiresAt) {
		metrics.Default.RecordCacheMiss()
		return nil, false
	}
	metrics.Default.RecordCacheHit()
	return e.value, true
}

func (m *Memory) gc() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		for k, e := range m.store {
			if time.Now().After(e.expiresAt) {
				delete(m.store, k)
			}
		}
		m.mu.Unlock()
	}
}
