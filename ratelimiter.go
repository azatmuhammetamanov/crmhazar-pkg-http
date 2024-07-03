package crmhazar_pkg_http

import (
	"sync"
	"time"
)

const (
	maxRequests     = 5
	rateLimitPeriod = time.Second
)

type RateLimiter struct {
	mu         sync.Mutex
	requests   map[string]int
	timestamps map[string]time.Time
	limit      int
	interval   time.Duration
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		requests:   make(map[string]int),
		timestamps: make(map[string]time.Time),
		limit:      maxRequests,
		interval:   rateLimitPeriod,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.interval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, ts := range rl.timestamps {
			if now.Sub(ts) > rl.interval {
				delete(rl.requests, key)
				delete(rl.timestamps, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(deviceID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if ts, found := rl.timestamps[deviceID]; found {
		if now.Sub(ts) > rl.interval {
			rl.requests[deviceID] = 0
			rl.timestamps[deviceID] = now
		}
	} else {
		rl.timestamps[deviceID] = now
	}

	if rl.requests[deviceID] < rl.limit {
		rl.requests[deviceID]++
		return true
	}

	return false
}
