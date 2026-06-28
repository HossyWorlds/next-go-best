---
name: api-contract-sync
description: GoのAPI(handler DTO・model)とTypeScriptのミラー型(frontend/lib/types)の整合をチェックし、ドリフトを検出するスキル。「型がずれてる」「API契約の同期」「フロントの型を合わせる」「contract sync」などで使用。
tools: Read, Glob, Grep, Bash, Edit
---

# api-contract-sync

`next-go-best` は OpenAPI を使わず、Go の DTO を `frontend/lib/types` に**手動ミラー**している。
このスキルは Go 側とTS側の型の食い違い（ドリフト）を検出し、TS側を更新する。

## チェック手順

### 1. 対応関係を把握

| Go | TypeScript |
|---|---|
| `backend/internal/<domain>/model.go` の JSON タグ付き struct | `frontend/lib/types/<domain>.ts` |
| `backend/internal/<domain>/handler.go` の request DTO | フォーム/Action の入力型 |
| エラー封筒 `{error:{code,message}}` | `lib/api/client.ts` `ApiErrorBody` |

### 2. Go 側の真実を抽出

各 model の **JSONフィールド名と型**を確認する（Goは `json:"camelCase"` を使用）。例:

```bash
# JSONタグ付きフィールドを一覧
grep -nE 'json:"' backend/internal/<domain>/model.go
```

注意点:
- `json:"-"` のフィールド（例: `PasswordHash`）は **TS型に含めない**。
- `time.Time` → TS では `string`（ISO8601）。
- `uuid.UUID` → `string`。
- enum（例: Status）→ TS の union 型と値域を一致させる。
- ネスト/一覧レスポンス（例: `ListResult{Tasks, Total}`）の形も合わせる。

### 3. ドリフトを照合

TS の型定義（`lib/types/<domain>.ts`）と突き合わせ、次を検出:
- フィールドの過不足（Go に有るが TS に無い／その逆）
- 名前の不一致（camelCase ずれ）
- 型の不一致（enum値・nullable・配列）
- `json:"-"` フィールドが TS に漏れていないか（**セキュリティ観点で重要**）

### 4. 修正

ドリフトがあれば **TS側を Go に合わせて** 更新する（Goが契約の真実）。
変更後 `cd frontend && pnpm lint && pnpm build` で型エラーが無いことを確認。

## 出力フォーマット

```
## API契約チェック結果: <domain>
- ✅ 一致 / ⚠️ ドリフト
- 検出: <Goフィールド> ↔ <TSフィールド> の差分
- 対応: TS側を…に修正（または要確認）
```

## 将来の改善提案（実行はしない）

手動ミラーは規模が増えると破綻しやすい。次フェーズで **OpenAPI 化**（Go側で spec 生成 → TS型を `openapi-typescript` で自動生成）を提案できる。必要なら別途相談。
