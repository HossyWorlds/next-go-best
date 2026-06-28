import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import type { AuthState } from "@/app/(auth)/actions";

import { AuthForm } from "./auth-form";

const noopAction = async (): Promise<AuthState> => ({});

describe("AuthForm", () => {
  it("メール・パスワード入力と送信ボタンを表示する", () => {
    render(
      <AuthForm
        title="ログイン"
        submitLabel="ログイン"
        action={noopAction}
        altHref="/register"
        altLabel="登録"
      />,
    );
    expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
    expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "ログイン" })).toBeInTheDocument();
  });
});
