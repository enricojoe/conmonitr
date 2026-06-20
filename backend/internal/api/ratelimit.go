package api

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	loginRate  = rate.Limit(5.0 / 60) // 5 attempts per minute
	loginBurst = 5
	ipTTL      = 10 * time.Minute
)

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type loginLimiter struct {
	mu      sync.Mutex
	entries map[string]*ipEntry
}

func newLoginLimiter() *loginLimiter {
	l := &loginLimiter{entries: make(map[string]*ipEntry)}
	go l.cleanup()
	return l
}

func (l *loginLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[ip]
	if !ok {
		e = &ipEntry{limiter: rate.NewLimiter(loginRate, loginBurst)}
		l.entries[ip] = e
	}
	e.lastSeen = time.Now()
	return e.limiter.Allow()
}

// cleanup removes entries that haven't been seen for ipTTL to prevent unbounded growth.
func (l *loginLimiter) cleanup() {
	for range time.Tick(ipTTL) {
		l.mu.Lock()
		for ip, e := range l.entries {
			if time.Since(e.lastSeen) > ipTTL {
				delete(l.entries, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *loginLimiter) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !l.allow(ip) {
			w.Header().Set("Retry-After", "60")
			jsonError(w, "too many login attempts, try again later", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
