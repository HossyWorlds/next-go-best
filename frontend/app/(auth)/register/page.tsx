import { AuthForm } from "@/components/auth/auth-form";

import { registerAction } from "../actions";

export default function RegisterPage() {
  return (
    <main className="flex min-h-screen items-center justify-center p-6">
      <AuthForm
        title="アカウント作成"
        submitLabel="登録する"
        action={registerAction}
        altHref="/login"
        altLabel="既にアカウントをお持ちの方"
      />
    </main>
  );
}
