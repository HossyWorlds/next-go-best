# フロント実装パターンリファレンス

## ディレクトリ責務

```
app/<route>/page.tsx      Server Component（データ取得・認証確認）
app/<route>/actions.ts    Server Actions（"use server"・変更系・revalidatePath）
components/<feature>/      機能コンポーネント（client/presentational）
components/ui/             shadcn/ui（編集しない・add で追加）
lib/api/                   Go API クライアント（サーバ側のみ）
lib/types/                 Go DTO のミラー型
```

## Server Action（変更系）

```ts
"use server";
import { revalidatePath } from "next/cache";
import { ApiError } from "@/lib/api/client";

export type FormState = { error?: string; ok?: boolean };

export async function createXAction(_prev: FormState, fd: FormData): Promise<FormState> {
  try {
    await createX({ /* fd から組み立て */ });
  } catch (e) {
    return { error: e instanceof ApiError ? e.message : "作成に失敗しました" };
  }
  revalidatePath("/x");
  return { ok: true };
}
```

- `redirect()` は例外で制御を移すため try/catch の**外**で呼ぶ。
- 削除のように引数が FormData だけなら `(fd: FormData) => Promise<void>` 形式にして `<form action={deleteX}>` で直接使う。

## フォームのサブミット制御（effect 内 setState を避ける）

```tsx
const [pending, startTransition] = useTransition();
const [error, setError] = useState<string>();

function handleSubmit(fd: FormData) {
  startTransition(async () => {
    const r = await action({}, fd);
    if (r.error) setError(r.error);
    else { setError(undefined); setOpen(false); } // ダイアログを閉じる等
  });
}
// <form action={handleSubmit}>
```

`react-hooks/set-state-in-effect` 違反を避けるため、結果処理は effect ではなくイベント（transition）内で行う。

## Server Component（取得・認証）

```tsx
export default async function Page() {
  const me = await getMe();
  if (!me) redirect("/login");
  const data = await listX();
  return /* ... */;
}
```

## API クライアント（Cookie 中継）

`lib/api/client.ts` の `apiFetch` が現在リクエストの Cookie を Go へ転送し、`cache: "no-store"`。
ログイン時は Go の `Set-Cookie` からトークンを取り出し Next 側 Cookie に再発行（`lib/api/auth.ts` 参照）。

## shadcn/ui の注意（base-ui ベース）

- ボタンで Link 等をラップするには `asChild` ではなく `render` プロップ: `<Button render={<Link href="/x">…</Link>} />`。
- Dialog は `open`/`onOpenChange` で制御、`DialogTrigger render={<Button/>}`。
- 新規コンポーネントは `pnpm dlx shadcn@latest add <name>` で追加。
