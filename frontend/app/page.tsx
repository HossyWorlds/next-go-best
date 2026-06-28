import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <main className="mx-auto flex min-h-screen max-w-2xl flex-col justify-center gap-8 p-6">
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">next-go-best</h1>
        <p className="text-muted-foreground">
          Next.js（App Router）+ Go（標準 net/http）のベストプラクティス雛形。
          セッションCookie認証つきの Tasks/Todo を題材に、レイヤリング・型安全・
          観測性・TDD を一通り通したポートフォリオの型です。
        </p>
        <ul className="list-inside list-disc text-sm text-muted-foreground">
          <li>フロント: Next.js 16 / React 19 / Tailwind v4 / shadcn/ui</li>
          <li>バック: Go / sqlc / pgx / Postgres / golang-migrate</li>
          <li>認証: 自前セッションCookie（argon2id）</li>
        </ul>
      </div>
      <div className="flex gap-3">
        <Button render={<Link href="/login">ログイン</Link>} />
        <Button
          variant="outline"
          render={<Link href="/register">アカウント作成</Link>}
        />
      </div>
    </main>
  );
}
