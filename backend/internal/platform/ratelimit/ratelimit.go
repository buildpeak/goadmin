package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	burst    int
	cleanup  time.Duration
}

type visitor struct {
	tokens   int
	lastSeen time.Time
}

func New(rate, burst int) *Limiter {
	l := &Limiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  5 * time.Minute,
	}
	go l.cleanupLoop()
	return l
}

func (l *Limiter) cleanupLoop() {
	for {
		time.Sleep(l.cleanup)
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > l.cleanup {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *Limiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[ip]
	now := time.Now()
	if !exists {
		l.visitors[ip] = &visitor{tokens: l.burst - 1, lastSeen: now}
		return true
	}

	elapsed := now.Sub(v.lastSeen)
	tokensToAdd := int(elapsed.Seconds()) * l.rate
	v.tokens += tokensToAdd
	if v.tokens > l.burst {
		v.tokens = l.burst
	}
	v.lastSeen = now

	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

func extractIP(req *http.Request) string {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}
	if fwd := req.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = fwd
	}
	return ip
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		if !l.allow(ip) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
