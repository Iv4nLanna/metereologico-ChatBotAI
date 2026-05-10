package metrics

import (
	"math"
	"sync/atomic"
	"time"
)

// bucketBounds defines latency histogram upper bounds in ms. -1 means ≥500ms.
var bucketBounds = [5]int64{10, 50, 100, 500, -1}

// Store holds live counters and a latency histogram.
// All atomic.Int64 fields are safe for concurrent access without locks.
type Store struct {
	totalRequests  atomic.Int64
	totalErrors    atomic.Int64
	cacheHits      atomic.Int64
	cacheMisses    atomic.Int64
	latencyBuckets [5]atomic.Int64 // <10ms, <50ms, <100ms, <500ms, ≥500ms
	startTime      time.Time
}

// Default is the process-wide metrics store, initialized at package load time.
var Default = &Store{startTime: time.Now()}

// Snapshot is the JSON-serializable view of current metrics.
type Snapshot struct {
	UptimeSeconds   int64   `json:"uptime_seconds"`
	TotalRequests   int64   `json:"total_requests"`
	TotalErrors     int64   `json:"total_errors"`
	ErrorRatePct    float64 `json:"error_rate_pct"`
	P95LatencyMs    int64   `json:"p95_latency_ms"`
	CacheHits       int64   `json:"cache_hits"`
	CacheMisses     int64   `json:"cache_misses"`
	CacheHitRatePct float64 `json:"cache_hit_rate_pct"`
}

// RecordRequest records a completed HTTP request. The counter and bucket
// increments are two separate atomic ops; a concurrent Snapshot() may
// observe a transient off-by-one, which is acceptable for a dashboard.
func (s *Store) RecordRequest(statusCode int, latencyMs int64) {
	s.totalRequests.Add(1)
	if statusCode >= 400 {
		s.totalErrors.Add(1)
	}
	for i, bound := range bucketBounds {
		if bound == -1 || latencyMs < bound {
			s.latencyBuckets[i].Add(1)
			break
		}
	}
}

func (s *Store) RecordCacheHit()  { s.cacheHits.Add(1) }
func (s *Store) RecordCacheMiss() { s.cacheMisses.Add(1) }

func (s *Store) Snapshot() Snapshot {
	total := s.totalRequests.Load()
	errors := s.totalErrors.Load()
	hits := s.cacheHits.Load()
	misses := s.cacheMisses.Load()

	var errorRate, hitRate float64
	if total > 0 {
		errorRate = round1(float64(errors) / float64(total) * 100)
	}
	if lookups := hits + misses; lookups > 0 {
		hitRate = round1(float64(hits) / float64(lookups) * 100)
	}

	return Snapshot{
		UptimeSeconds:   int64(time.Since(s.startTime).Seconds()),
		TotalRequests:   total,
		TotalErrors:     errors,
		ErrorRatePct:    errorRate,
		P95LatencyMs:    s.p95(),
		CacheHits:       hits,
		CacheMisses:     misses,
		CacheHitRatePct: hitRate,
	}
}

func (s *Store) p95() int64 {
	var counts [5]int64
	var total int64
	for i := range s.latencyBuckets {
		counts[i] = s.latencyBuckets[i].Load()
		total += counts[i]
	}
	if total == 0 {
		return 0
	}
	threshold := int64(math.Ceil(float64(total) * 0.95))
	cumulative := int64(0)
	bounds := [5]int64{10, 50, 100, 500, 500}
	for i, c := range counts {
		cumulative += c
		if cumulative >= threshold {
			return bounds[i]
		}
	}
	return 500
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
