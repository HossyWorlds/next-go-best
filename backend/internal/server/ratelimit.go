package server

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter はクライアントIP単位の固定ウィンドウ・レート制限。
//
// 注意: これはプロセス内（単一インスタンス）の簡易実装。
// 複数インスタンス構成では Redis 等の共有ストアに置き換えること。
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	max      int
	window   time.Duration
}

type visitor struct {
	count   int
	resetAt time.Time
}

func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      max,
		window:   window,
	}
}

func (rl *RateLimiter) allow(key string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[key]
	if !ok || now.After(v.resetAt) {
		rl.visitors[key] = &visitor{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	if v.count >= rl.max {
		return false
	}
	v.count++
	return true
}

// Middleware はクライアントIPでレート制限する。超過時は 429 を返す。
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r), time.Now()) {
			writeErrorEnvelope(w, http.StatusTooManyRequests, "rate_limited", "リクエストが多すぎます。しばらく待って再試行してください")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
