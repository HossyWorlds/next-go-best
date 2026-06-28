package server

import (
	"context"
	"net/http"
	"time"
)

// Pinger は readiness 判定のための疎通確認（*pgxpool.Pool が満たす）。
type Pinger interface {
	Ping(ctx context.Context) error
}

// Healthz は liveness。プロセスが生きていれば 200。
func Healthz(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz は readiness。依存（DB）が応答可能なら 200、不可なら 503。
func Readyz(p Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := p.Ping(ctx); err != nil {
			writeErrorEnvelope(w, http.StatusServiceUnavailable, "not_ready", "依存サービスに接続できません")
			return
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}
