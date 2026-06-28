# next-go-best

**Next.js（フロント）+ Go（バックエンド）のベストプラクティス雛形** — ポートフォリオの「型」。

公式ドキュメント・有名OSSの作法に沿った、実務水準で語れる最小構成。題材は汎用CRUDの **Tasks/Todo** ＋ **セッションCookie認証**。開発は **TDD**（Red→Green→Refactor）前提。

---

## 技術スタック

| レイヤ | 採用 |
|---|---|
| フロント | Next.js 16 (App Router) / React 19 / TypeScript strict / Tailwind v4 / shadcn/ui |
| バックエンド | Go (標準 `net/http` ServeMux) / layered構成 / `log/slog` |
| データ | Postgres 16 / sqlc / golang-migrate / pgx v5 |
| 認証 | 自前セッションCookie（argon2id + サーバサイドセッション） |
| テスト | Go `testing`+`httptest`+testcontainers-go / Vitest+Testing Library / Playwright |

## クイックスタート

```bash
cp .env.example .env          # ローカルはそのままでOK
make up                       # Postgres 起動（docker）
make migrate                  # マイグレーション適用
make sqlc                     # クエリからGoコード生成（初回/クエリ変更時）
make dev-api                  # API: http://localhost:8080

# 別ターミナル
make front-install
make dev-front                # Web: http://localhost:3000
```

> ローカルに別の Postgres が 5432 で動いている場合は `.env` の `POSTGRES_PORT` と `DATABASE_URL` を 5433 等に変更する。

## ディレクトリ構成

```
backend/
  cmd/api/            APIエントリ（DI配線 + graceful shutdown）
  cmd/migrate/        埋め込みSQLでのマイグレーション実行
  internal/
    app/              ルート配線（main・統合テストが共用）+ 統合テスト
    server/           ミドルウェア・レスポンス封筒・health・rate limit
    auth/             認証スライス（argon2id / セッション / middleware）
    task/             タスクスライス（handler/service/repository/sqlc実装）
    apperr/           ドメインエラー（→HTTPステータス変換は server 層）
    validate/         go-playground/validator ラッパ
    db/               接続プール・マイグレーション・sqlc生成コード(gen)
  migrations/         golang-migrate（*.up.sql / *.down.sql）
frontend/
  app/                App Router（page=取得 / actions=変更）
    (auth)/           login・register
    tasks/            一覧ページ + Server Actions
  components/         tasks/ ・ auth/ ・ ui(shadcn)
  lib/api/            Go API クライアント（サーバ側・Cookie中継）
  lib/types/          Go DTO のミラー型
  e2e/                Playwright
  proxy.ts            Next16の保護ルート（旧 middleware.ts）
.claude/skills/       使い回し用 Claude Code skills（後述）
```

## 設計の要点

### バックエンド
- **レイヤリング**: handler（HTTP境界）→ service（ロジック・所有者チェック・検証）→ repository(interface) → postgres実装(sqlc)。service は interface のみに依存し fake でテスト可能。
- **エラー**: ドメインは `apperr`（Validation/Unauthorized/Forbidden/NotFound/Conflict/Internal）を返し、`server.WriteError` が HTTPステータス＋封筒 `{"error":{"code","message"}}` に一元変換。
- **観測性**: `log/slog` 構造化ログ・requestID ミドルウェア・`/healthz`(liveness)・`/readyz`(DB ping)。
- **セキュリティ**: CORS は許可オリジン限定（ワイルドカード禁止）。一覧は必ず `LIMIT/OFFSET`。秘密情報は `.env`（ダミーのみコミット）。
- **認証**: パスワードは argon2id。セッションはランダムトークンを発行し、DBには**ハッシュ**を保存。Cookie は `HttpOnly`+`SameSite=Lax`+（本番）`Secure`。ログインにレート制限・パスワードポリシー・セッション失効。Task は全て要認証＋所有者チェック。

### フロントエンド（Next.js 16 の作法）
- データ取得は **Server Component**、変更は **Server Actions** ＋ `revalidatePath`。
- `proxy.ts`（旧 `middleware.ts`）で保護ルートを制御（nodejsランタイム）。
- `cookies()` 等は **async**。`next lint` は廃止 → `pnpm lint`（ESLint CLI）。Turbopack デフォルト。
- API はサーバ側（`lib/api`）から呼び、セッションCookieを Go へ中継。トークンをブラウザに露出しない。

## エンドポイント

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/healthz` / `/readyz` | liveness / readiness |
| POST | `/api/v1/auth/register` | 登録 |
| POST | `/api/v1/auth/login` | ログイン（Cookie発行・レート制限） |
| POST | `/api/v1/auth/logout` | ログアウト |
| GET | `/api/v1/auth/me` | 現在のユーザー（要認証） |
| GET | `/api/v1/tasks` | 一覧（自分のみ・要認証） |
| POST | `/api/v1/tasks` | 作成（要認証） |
| GET/PUT/DELETE | `/api/v1/tasks/{id}` | 詳細/更新/削除（所有者のみ） |

## TDD フロー

各スライスを **Red → Green → Refactor** で実装する。

1. **単体テストを書く（Red）**: service を fake repository で（`task/service_test.go`）。
2. **統合テストを書く（Red）**: testcontainers で実Postgresを起動し `httptest` でフル経路（`app/*_integration_test.go`）。
3. **実装（Green）**: handler/service/repository/sqlc を作り、テストを通す。
4. **Refactor**: 規約（命名・責務・エラー）を整える。

フロントも同様に Vitest（コンポーネント）→ Playwright（フロー）を先に書く。

## テスト

```bash
make test-back                  # Go 単体 + 統合(testcontainers・Docker必須)
cd backend && go test -short ./...   # 統合を除く（Docker無し環境）

make test-front                 # Vitest（コンポーネント）
cd frontend && pnpm exec playwright install chromium  # 初回のみ
make e2e                        # Playwright（要: Postgres + API 起動）

make lint-back                  # golangci-lint（go tool）
make lint-front                 # ESLint
make build-front                # 本番ビルド
```

## Claude Code Skills

`.claude/skills/` に、今後の開発で使い回せるスキルを同梱（日本語）。

| skill | 用途 |
|---|---|
| `scaffold-go-endpoint` | 新ドメインの垂直スライスを規約どおり **TDD順** で生成 |
| `scaffold-next-feature` | App Router の新機能（page/actions/components/型）を **テスト先行** で生成 |
| `api-contract-sync` | Go DTO と `lib/types` の TS 型のドリフトを検出・修正 |
| `go-next-review` | この雛形の規約（命名/レイヤ/エラー/CORS/LIMIT/観測性/Next16作法）でレビュー |

使い方: Claude Code で `/scaffold-go-endpoint` のように呼ぶ（または自然文で依頼すると description にマッチして起動）。

## 拡張ポイント（この雛形では範囲外）

- メール検証・パスワードリセット・CSRFダブルサブミット
- マルチインスタンス向けの分散レート制限（Redis 等）
- 本番デプロイ（Dockerfile / CI(GitHub Actions) / IAM・シークレット管理・TLS）
- API契約の OpenAPI 化（手動ミラー型の自動生成）
# next-go-best
