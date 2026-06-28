---
name: go-next-review
description: next-go-best の規約(命名/レイヤ責務/エラーハンドリング/CORS/LIMIT/観測性/Next16作法)に沿ってコードをレビューするスキル。「レビューして」「規約チェック」「この実装大丈夫?」「go-next-review」などで使用。バグ探しではなく規約準拠の確認に特化する。
tools: Read, Glob, Grep, Bash
---

# go-next-review

`next-go-best` の規約準拠をチェックするレビュー用スキル。
（一般的なバグ探しは `/code-review` を使う。こちらは**この雛形固有の規約**に特化）

## 進め方

1. レビュー対象を特定（`git diff` または指定ファイル）。
2. 下のチェックリストを Backend/Frontend で適用。
3. 違反は「該当箇所・規約・推奨修正」をセットで指摘。事実と推測を分ける。

## Backend チェックリスト

**レイヤリング**
- [ ] handler はHTTP境界のみ（ロジックを持たない）
- [ ] service は `repository` interface にのみ依存（pgx/sql を直接触らない）
- [ ] repository_postgres は sqlc 生成を domain 型へ変換する薄い層

**エラー / セキュリティ**
- [ ] ドメイン層は `apperr` を返す（HTTPステータスを直接書かない）
- [ ] HTTP変換は `server.WriteError` 一元化
- [ ] 5xx 以外でスタックや機密をレスポンスに出していない（PII/トークン/パスワード）
- [ ] 所有者スコープのリソースは service で owner チェック（403）
- [ ] CORS は許可オリジンのみ（ワイルドカード禁止）
- [ ] 認証必須ルートは `auth.Middleware` 配下
- [ ] パスワードは argon2id、セッションはトークンの**ハッシュ**保存

**DB / 観測性**
- [ ] 一覧クエリに `LIMIT/OFFSET` デフォルト（全件返却なし）
- [ ] migration に up/down 両方、ゼロパディング連番
- [ ] `log/slog` 構造化ログ＋requestID、healthz/readyz がある

**テスト**
- [ ] service 単体（fake repo）＋統合（testcontainers）が両方ある
- [ ] 403/404/400 の異常系を網羅

## Frontend チェックリスト

**Next.js 16 作法**
- [ ] `middleware.ts` ではなく `proxy.ts`（関数名 `proxy`）
- [ ] `cookies()/headers()/params/searchParams` を `await` している
- [ ] `next lint` ではなく ESLint CLI（`pnpm lint`）

**設計**
- [ ] データ取得は Server Component、変更は Server Action＋`revalidatePath`
- [ ] API 呼び出しはサーバ側（`lib/api`）でCookie中継。トークンをブラウザに露出しない
- [ ] フォームの成功/失敗処理を effect 内 setState で行っていない（`useTransition`）
- [ ] 保護ルートは `proxy.ts` の matcher に登録
- [ ] `lib/types` が Go DTO とずれていない（`api-contract-sync` 参照）

**型 / 品質**
- [ ] TypeScript strict を満たす（`pnpm build` が通る）
- [ ] `json:"-"` の項目（password等）が TS 型・画面に漏れていない

## 実行コマンド（裏取り）

```bash
cd backend && go vet ./... && go test ./...      # 統合はDocker必要。なければ -short
cd frontend && pnpm lint && pnpm test && pnpm build
```

## 出力フォーマット

```
## go-next-review 結果
### 🔴 規約違反
- <file:line> 規約「…」に違反。推奨: …
### 🟡 改善提案
- …
### ✅ 良い点
- …
```
