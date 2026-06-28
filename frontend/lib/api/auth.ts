import { cookies } from "next/headers";

import { apiFetch, apiJSON, toApiError, SESSION_COOKIE } from "./client";
import type { User } from "@/lib/types/user";

// 認証API。login/logout は Server Action からのみ呼ぶこと（Cookie書き込みのため）。

export function register(email: string, password: string): Promise<User> {
  return apiJSON<User>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function login(email: string, password: string): Promise<User> {
  const res = await apiFetch("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    throw await toApiError(res);
  }
  const user = (await res.json()) as User;

  // Go が発行した Set-Cookie からトークンを取り出し、Next オリジンのCookieとして再発行する。
  const token = parseSessionToken(res.headers.get("set-cookie"));
  if (token) {
    const store = await cookies();
    store.set(SESSION_COOKIE, token, {
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      path: "/",
    });
  }
  return user;
}

export async function logout(): Promise<void> {
  // Go 側のセッションも失効させる（Cookieは現在のリクエストから転送される）。
  await apiFetch("/api/v1/auth/logout", { method: "POST" });
  const store = await cookies();
  store.delete(SESSION_COOKIE);
}

export async function getMe(): Promise<User | null> {
  const res = await apiFetch("/api/v1/auth/me");
  if (res.status === 401) {
    return null;
  }
  if (!res.ok) {
    throw await toApiError(res);
  }
  return (await res.json()) as User;
}

function parseSessionToken(setCookie: string | null): string | null {
  if (!setCookie) {
    return null;
  }
  const match = setCookie.match(/(?:^|,\s*)session=([^;]+)/);
  return match ? match[1] : null;
}
