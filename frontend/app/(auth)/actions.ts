"use server";

import { redirect } from "next/navigation";

import { login, register } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";

export type AuthState = { error?: string };

// loginAction はログインし、成功したら /tasks へ遷移する。
export async function loginAction(
  _prev: AuthState,
  formData: FormData,
): Promise<AuthState> {
  const email = String(formData.get("email") ?? "");
  const password = String(formData.get("password") ?? "");

  try {
    await login(email, password);
  } catch (e) {
    return { error: errorMessage(e, "ログインに失敗しました") };
  }
  // redirect は例外で制御を移すため try/catch の外で呼ぶ。
  redirect("/tasks");
}

// registerAction は登録後そのままログインし、/tasks へ遷移する。
export async function registerAction(
  _prev: AuthState,
  formData: FormData,
): Promise<AuthState> {
  const email = String(formData.get("email") ?? "");
  const password = String(formData.get("password") ?? "");

  try {
    await register(email, password);
    await login(email, password);
  } catch (e) {
    return { error: errorMessage(e, "登録に失敗しました") };
  }
  redirect("/tasks");
}

function errorMessage(e: unknown, fallback: string): string {
  return e instanceof ApiError ? e.message : fallback;
}
