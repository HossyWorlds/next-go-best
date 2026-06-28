---
name: scaffold-next-feature
description: Next.js(App Router)に新しい機能(ページ/コンポーネント/Server Actions/型/APIクライアント)をTDDで追加するスキル。「画面を追加」「ページを作る」「新機能のUI」「scaffold next」などで使用。next-go-best のフロント規約(Server Components + Server Actions + shadcn/ui)に沿って生成する。
tools: Read, Glob, Grep, Bash, Edit, Write
---

# scaffold-next-feature

`next-go-best` のフロント規約に沿って、新しい機能を **テスト先行** で追加する。

## 前提の確認

1. `frontend/app/tasks/` と `frontend/components/tasks/` を**参照実装**として読む。
2. 機能名・ルート・必要な画面（一覧/作成/編集など）を確認する。
3. 対応する Go API があるか（無ければ先に `scaffold-go-endpoint`）。

## Next.js 16 の作法（必守）

- `middleware.ts` は使わない → `proxy.ts`（関数名 `proxy`、nodejsランタイム）。保護ルートは matcher に追加。
- `next lint` は廃止 → `pnpm lint`（ESLint CLI）。
- `cookies()`/`headers()`/`params`/`searchParams` は **async**（`await` 必須）。
- データ取得は Server Component、変更は Server Actions（`"use server"`）＋ `revalidatePath`。
- 重いクライアント状態ライブラリは入れない。

## 生成手順（この順で）

### 1. 型（Go DTOのミラー）

`frontend/lib/types/<feature>.ts` に Go のレスポンス型を写す。`lib/types/task.ts` を手本に。
（同期は `api-contract-sync` skill で検証）

### 2. APIクライアント

`frontend/lib/api/<feature>.ts` に `apiJSON`/`apiFetch`（`lib/api/client.ts`）を使ったラッパを作る。サーバ側専用・Cookie中継。

### 3. テストを先に書く（Red）

- `frontend/components/<feature>/*.test.tsx`: Vitest + Testing Library。`components/tasks/*.test.tsx` を手本に、描画・初期値・ラベルを検証。
- `frontend/e2e/<feature>.spec.ts`: Playwright でコアフロー（必要なら認証込み）。`e2e/tasks.spec.ts` を手本に。

`pnpm test` が失敗することを確認。

### 4. 実装（Green）

- Server Actions: `app/<route>/actions.ts`（`"use server"`、`createTaskAction` を手本に。成功時 `revalidatePath`）。
- コンポーネント: `components/<feature>/`。shadcn/ui（`components/ui`）を使う。フォームは `useTransition` でサブミット制御（effect内 setState を避ける）。
- ページ: `app/<route>/page.tsx`（Server Component。`getMe()` で認証確認 → 未認証は `redirect("/login")`）。
- 保護ルートなら `proxy.ts` の matcher に追加。

### 5. Green 確認

`pnpm test`（Vitest）→ `pnpm lint` → `pnpm build` を通す。e2e は API+DB 起動のうえ `pnpm exec playwright test`。

## 規約チェック

- [ ] データ取得は Server Component、変更は Server Action＋`revalidatePath`
- [ ] `cookies()` 等は `await` している
- [ ] フォームの成功/失敗が表示され、effect 内 setState をしていない
- [ ] 保護ルートは `proxy.ts` に登録
- [ ] Vitest（コンポーネント）と Playwright（フロー）の両方がある
- [ ] `pnpm lint` と `pnpm build` が通る

UI/状態の具体パターンは `references/patterns.md` を参照。
