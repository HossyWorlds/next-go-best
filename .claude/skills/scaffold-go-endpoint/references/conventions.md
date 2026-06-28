# バックエンド規約リファレンス

## レイヤと依存方向

```
handler ──> service ──> repository(interface) <── repository_postgres (sqlc)
   │           │
   └─ HTTP境界  └─ ビジネスロジック・所有者チェック・バリデーション
```

- handler は HTTP の詳細（デコード/ステータス/Cookie）のみ。ロジックを持たない。
- service は `repository` interface にのみ依存（DB非依存）。fake でテスト可能に保つ。
- repository_postgres は sqlc 生成コードを domain 型へ変換する薄い層。

## エラーハンドリング

- ドメイン層は `internal/apperr` を返す（`Validation/Unauthorized/Forbidden/NotFound/Conflict/Internal`）。
- HTTPステータスへの変換は `server.WriteError` が一元化（`statusFromKind`）。
- レスポンス封筒は `{ "error": { "code", "message" } }`。
- 5xx は原因を `slog` に出すが、レスポンス本文・ログに機密（PII/トークン）を出さない。

## 所有者スコープ（owner_id）

- repository は ID のみで取得（`GetByID`）。
- service が `row.OwnerID != userID` を判定し `apperr.Forbidden`（403）を返す。
- 存在しなければ `apperr.NotFound`（404）。
- 一覧は `WHERE owner_id = $1` で必ず絞る。

## ページネーション

- service に `DefaultLimit=20` / `MaxLimit=100`。`normalizeLimit` で丸める。
- 一覧クエリは `ORDER BY created_at DESC LIMIT $2 OFFSET $3`。

## バリデーション

- 入力 struct に `validate` タグ（go-playground/validator）。
- service 冒頭で `validate.Struct(p)`。失敗は自動で `apperr.Validation`。

## テスト

- 単体: `service_test.go` に fakeRepo（インメモリ）＋ table-driven。
- 統合: `internal/app/<domain>_integration_test.go`。testcontainers で実Postgres、`httptest` でフル経路。
  - `-short` で統合をスキップ（`TestMain` 参照）。
- 必ず「所有者外アクセス403」「存在しないID 404」「バリデーション400」を網羅。

## sqlc / migration の注意

- migration は `*.up.sql` / `*.down.sql`（golang-migrate形式）。連番はゼロパディング。
- sqlc は down を無視し、schema を `migrations` から読む（`backend/sqlc.yaml`）。
- uuid→`google/uuid.UUID`、timestamptz→`time.Time` に override 済み。
