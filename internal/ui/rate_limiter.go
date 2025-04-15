package ui

import (
	"net/url"
	"sync"
	"time"
)

// RateLimiter implements a per-domain rate limiter
type RateLimiter struct {
	mu            sync.Mutex
	requestsPerDomain map[string][]time.Time
	requestsPerSecond int
	windowSeconds     int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		requestsPerDomain: make(map[string][]time.Time),
		requestsPerSecond: requestsPerSecond,
		windowSeconds:     windowSeconds,
	}
}

// Allow checks if a request is allowed for the given URL
func (rl *RateLimiter) Allow(urlStr string) bool {
	// Parse the URL to get the domain
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// If we can't parse the URL, allow the request
		return true
	}
	
	domain := parsedURL.Host
	
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Get the current time
	now := time.Now()
	
	// Clean up old requests
	rl.cleanupOldRequests(domain, now)
	
	// Check if we've exceeded the rate limit
	if len(rl.requestsPerDomain[domain]) >= rl.requestsPerSecond {
		return false
	}
	
	// Add the current request
	rl.requestsPerDomain[domain] = append(rl.requestsPerDomain[domain], now)
	
	return true
}

// cleanupOldRequests removes requests older than the window
func (rl *RateLimiter) cleanupOldRequests(domain string, now time.Time) {
	cutoff := now.Add(-time.Duration(rl.windowSeconds) * time.Second)
	
	requests, ok := rl.requestsPerDomain[domain]
	if !ok {
		return
	}
	
	// Find the index of the first request that's within the window
	i := 0
	for ; i < len(requests); i++ {
		if requests[i].After(cutoff) {
			break
		}
	}
	
	// Remove all requests before the index
	if i > 0 {
		rl.requestsPerDomain[domain] = requests[i:]
	}
}
