import { AuthForm } from "@/components/auth/auth-form";

import { loginAction } from "../actions";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center p-6">
      <AuthForm
        title="ログイン"
        submitLabel="ログイン"
        action={loginAction}
        altHref="/register"
        altLabel="アカウントを作成する"
      />
    </main>
  );
}
