// Package config は環境変数からアプリ設定を読み込む。
// 秘密情報はコードに埋め込まず、すべて環境変数（.env）から渡す。
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config はアプリ全体の設定。
type Config struct {
	Port           string        // APIサーバの待受ポート
	DatabaseURL    string        // Postgres接続文字列
	AllowedOrigins []string      // CORSで許可するオリジン（ワイルドカード禁止）
	CookieSecure   bool          // セッションCookieの Secure 属性（本番は true）
	SessionTTL     time.Duration // セッション有効期限
}

// Load は環境変数から Config を構築する。必須値が欠けていればエラーを返す。
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL is required")
	}

	ttlHours := getInt("SESSION_TTL_HOURS", 72)

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    dbURL,
		AllowedOrigins: splitAndTrim(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		CookieSecure:   getBool("SESSION_COOKIE_SECURE", false),
		SessionTTL:     time.Duration(ttlHours) * time.Hour,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
