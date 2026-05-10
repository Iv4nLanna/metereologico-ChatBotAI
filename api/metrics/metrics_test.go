package metrics

import (
	"testing"
	"time"
)

func TestRecordRequest_countsTotal(t *testing.T) {
	s := &Store{startTime: time.Now()}
	s.RecordRequest(200, 5)
	s.RecordRequest(200, 20)
	snap := s.Snapshot()
	if snap.TotalRequests != 2 {
		t.Fatalf("want 2 total requests, got %d", snap.TotalRequests)
	}
}

func TestRecordRequest_countsErrors(t *testing.T) {
	s := &Store{startTime: time.Now()}
	s.RecordRequest(200, 5)
	s.RecordRequest(400, 5)
	s.RecordRequest(500, 5)
	snap := s.Snapshot()
	if snap.TotalErrors != 2 {
		t.Fatalf("want 2 errors, got %d", snap.TotalErrors)
	}
	if snap.ErrorRatePct != 66.7 {
		t.Fatalf("want 66.7%% error rate, got %.1f%%", snap.ErrorRatePct)
	}
}

func TestCacheHitRate(t *testing.T) {
	s := &Store{startTime: time.Now()}
	s.RecordCacheHit()
	s.RecordCacheHit()
	s.RecordCacheHit()
	s.RecordCacheMiss()
	snap := s.Snapshot()
	if snap.CacheHits != 3 {
		t.Fatalf("want 3 hits, got %d", snap.CacheHits)
	}
	if snap.CacheHitRatePct != 75.0 {
		t.Fatalf("want 75.0%% hit rate, got %.1f%%", snap.CacheHitRatePct)
	}
}

func TestP95_allFast(t *testing.T) {
	s := &Store{startTime: time.Now()}
	for i := 0; i < 100; i++ {
		s.RecordRequest(200, 5) // all < 10ms → bucket[0]
	}
	snap := s.Snapshot()
	if snap.P95LatencyMs != 10 {
		t.Fatalf("want p95=10ms, got %d", snap.P95LatencyMs)
	}
}

func TestP95_slowTail(t *testing.T) {
	s := &Store{startTime: time.Now()}
	for i := 0; i < 94; i++ {
		s.RecordRequest(200, 5) // < 10ms → bucket[0]
	}
	for i := 0; i < 6; i++ {
		s.RecordRequest(200, 600) // ≥ 500ms → bucket[4]
	}
	// threshold = ceil(100 * 0.95) = 95; bucket[0] cumulative = 94 < 95 → falls to bucket[4]
	snap := s.Snapshot()
	if snap.P95LatencyMs != 500 {
		t.Fatalf("want p95=500ms, got %d", snap.P95LatencyMs)
	}
}

func TestP95_noRequests(t *testing.T) {
	s := &Store{startTime: time.Now()}
	snap := s.Snapshot()
	if snap.P95LatencyMs != 0 {
		t.Fatalf("want p95=0 when no requests, got %d", snap.P95LatencyMs)
	}
}
