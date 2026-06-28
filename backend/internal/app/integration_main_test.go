package app_test

// 統合テスト: testcontainers で使い捨ての実Postgresを起動し、マイグレーション適用後に
// 実際に配線した http.Handler（auth + task）を httptest 経由で叩く。
// HTTPハンドラ→service→repository→DB の全経路を検証する。
//
// Docker が必要。Docker が無い環境では `-short` でスキップする。

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/HossyWorlds/next-go-best/backend/internal/app"
	"github.com/HossyWorlds/next-go-best/backend/internal/config"
	"github.com/HossyWorlds/next-go-best/backend/internal/db"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	flag.Parse() // testing.Short() を使う前にフラグを解釈する
	if testing.Short() {
		// -short では統合テストをスキップ（Docker不要のCI段階用）。
		os.Exit(0)
	}
	os.Exit(runWithPostgres(m))
}

func runWithPostgres(m *testing.M) int {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("app"),
		tcpostgres.WithUsername("app"),
		tcpostgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Printf("failed to start postgres container (Docker required): %v", err)
		return 1
	}
	defer func() { _ = testcontainers.TerminateContainer(container) }()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("connection string: %v", err)
		return 1
	}

	if err := db.Migrate(dsn); err != nil {
		log.Printf("migrate: %v", err)
		return 1
	}

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Printf("pool: %v", err)
		return 1
	}
	defer pool.Close()

	cfg := &config.Config{
		Port:           "0",
		DatabaseURL:    dsn,
		AllowedOrigins: []string{"http://localhost:3000"},
		CookieSecure:   false,
		SessionTTL:     72 * time.Hour,
	}
	handler := app.NewHandler(app.Deps{Pool: pool, Cfg: cfg, Logger: discardLogger()})

	testServer = httptest.NewServer(handler)
	defer testServer.Close()

	return m.Run()
}

// ---- テストヘルパー ----

// newClient は Cookie を保持するHTTPクライアントを返す（セッション継続用）。
func newClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar: %v", err)
	}
	return &http.Client{Jar: jar}
}

// doJSON は JSON リクエストを送り、レスポンスのステータスと本文を返す。
func doJSON(t *testing.T, c *http.Client, method, path string, body any) (int, []byte) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		reader = strings.NewReader(string(b))
	}
	req, err := http.NewRequest(method, testServer.URL+path, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, data
}

// registerAndLogin は新規ユーザーを作成しログインまで行う（Cookieはclientに保持される）。
func registerAndLogin(t *testing.T, c *http.Client, email, password string) {
	t.Helper()
	if status, body := doJSON(t, c, http.MethodPost, "/api/v1/auth/register",
		map[string]string{"email": email, "password": password}); status != http.StatusCreated {
		t.Fatalf("register: status=%d body=%s", status, body)
	}
	if status, body := doJSON(t, c, http.MethodPost, "/api/v1/auth/login",
		map[string]string{"email": email, "password": password}); status != http.StatusOK {
		t.Fatalf("login: status=%d body=%s", status, body)
	}
}

// uniqueEmail はテストごとに衝突しないメールアドレスを生成する。
func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano())
}

// discardLogger はテスト中のログ出力を捨てる。
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
