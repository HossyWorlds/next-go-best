// Package app は依存を配線して http.Handler を組み立てる。
// cmd/api（本番起動）と統合テストの双方が同じ配線を再利用する。
package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HossyWorlds/next-go-best/backend/internal/auth"
	"github.com/HossyWorlds/next-go-best/backend/internal/config"
	"github.com/HossyWorlds/next-go-best/backend/internal/server"
	"github.com/HossyWorlds/next-go-best/backend/internal/task"
)

// Deps は Handler 構築に必要な依存。
type Deps struct {
	Pool   *pgxpool.Pool
	Cfg    *config.Config
	Logger *slog.Logger
}

// NewHandler は全ルートとグローバルミドルウェアを組み立てて返す。
func NewHandler(d Deps) http.Handler {
	authSvc := auth.NewService(auth.NewPostgresRepository(d.Pool), d.Cfg.SessionTTL)
	authHandler := auth.NewHandler(authSvc, d.Cfg.CookieSecure)

	taskSvc := task.NewService(task.NewPostgresRepository(d.Pool))
	taskHandler := task.NewHandler(taskSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", server.Healthz)
	mux.Handle("GET /readyz", server.Readyz(d.Pool))

	// ログインは総当たり対策にレート制限（IP単位・10回/分）。
	loginLimiter := server.NewRateLimiter(10, time.Minute)
	authHandler.RegisterRoutes(mux, loginLimiter.Middleware, authSvc.Middleware)
	taskHandler.RegisterRoutes(mux, authSvc.Middleware)

	// グローバルミドルウェア（外側→内側）。
	return server.Chain(mux,
		server.RequestID,
		server.Logger(d.Logger),
		server.Recoverer(d.Logger),
		server.CORS(d.Cfg.AllowedOrigins),
	)
}
