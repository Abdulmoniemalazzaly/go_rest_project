package middlewares

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	go rl.resetVisitorCount()
	return rl
}

func (rl *RateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		visitorIp := r.RemoteAddr
		rl.visitors[visitorIp]++
		if rl.visitors[visitorIp] > rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
