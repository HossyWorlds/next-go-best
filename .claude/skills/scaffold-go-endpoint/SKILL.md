---
name: scaffold-go-endpoint
description: Goバックエンドに新しいRESTエンドポイント(ドメインの垂直スライス)をTDDで追加するスキル。「新しいエンドポイント」「APIを追加」「ドメインを追加」「scaffold go」などで使用。next-go-best のレイヤリング規約(handler/service/repository/sqlc)に沿って、テスト先行で生成する。
tools: Read, Glob, Grep, Bash, Edit, Write
---

# scaffold-go-endpoint

`next-go-best` のバックエンド規約に沿って、新しいドメインの垂直スライスを **TDD（Red→Green→Refactor）** で追加する。

## 前提の確認

1. `backend/internal/task/` を**参照実装**として必ず読む（命名・責務分割・エラーハンドリングの手本）。
2. 追加するドメイン名をユーザーに確認する（単数形 例: `project`、複数形 例: `projects`）。
3. 認証必須か・所有者スコープ（owner_id）が要るかを確認する。

## 生成手順（この順で進める）

### 1. マイグレーション（ロールバック必須）

`backend/migrations/` に次の連番でゼロパディング命名（sqlc は辞書順で読むため）:

- `NNNNNN_create_<plural>.up.sql` / `.down.sql`

既存の `000002_create_tasks.up.sql` を雛形にする。owner所有なら `owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE` と `idx_<plural>_owner_created` を付ける。

### 2. sqlc クエリ

`backend/internal/db/queries/<domain>.sql` を作成。`task.sql` を手本に、一覧は必ず `LIMIT/OFFSET` と owner 絞り込みを入れる。
その後 `cd backend && go tool sqlc generate` で `internal/db/gen` を再生成する。

### 3. テストを先に書く（Red）

- `backend/internal/<domain>/service_test.go`: fake repository で table-driven。`task/service_test.go` の構造（fakeRepo, kindOf ヘルパ）を流用。
- `backend/internal/app/<domain>_integration_test.go`: 実Postgres(testcontainers)経由のCRUD＋所有者外403。`app/task_integration_test.go` を手本に。

この時点で `go test ./...` が**失敗する**ことを確認する。

### 4. 実装（Green）

`task/` の各ファイルを対応させて作る:

| ファイル | 役割 |
|---|---|
| `model.go` | ドメインモデル＋CreateParams/UpdateParams（validateタグ） |
| `repository.go` | interface＋`ErrNotFound` |
| `repository_postgres.go` | sqlc実装。`pgx.ErrNoRows`→`ErrNotFound` |
| `service.go` | バリデーション・所有者チェック（403/404）・ページネーション |
| `handler.go` | DecodeJSON→service→WriteJSON。`auth.RequireUserID` で認証ユーザー取得 |

エラーは必ず `apperr`（Validation/NotFound/Forbidden/Conflict/Internal）で返し、HTTP変換は `server.WriteError` に委ねる。

### 5. ルート配線

`backend/internal/app/app.go` の `NewHandler` に `<domain>Handler.RegisterRoutes(mux, authSvc.Middleware)` を追加。

### 6. Green 確認

`make up && make migrate && cd backend && go test ./...` が全緑になるまで直す。

## 規約チェック（最後に確認）

- [ ] 一覧に `LIMIT/OFFSET` デフォルトがある（全件返却しない）
- [ ] 認証必須ルートは `auth.Middleware` 配下
- [ ] 所有者スコープのリソースは service で owner チェック（403）
- [ ] migration に down がある
- [ ] service 単体テスト＋統合テストの両方がある
- [ ] エラーは `apperr` 経由、ログに機密を出さない

詳細な対応表は `references/conventions.md` を参照。
