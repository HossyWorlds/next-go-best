import { cookies } from "next/headers";

// サーバ側からのみ呼ぶ Go API クライアント。
// ブラウザ → Next(3000) → Go(8080) の経路で、セッションCookieを中継する。
const API_URL = process.env.API_URL ?? "http://localhost:8080";

export const SESSION_COOKIE = "session";

// Go の共通エラー封筒: { "error": { "code", "message" } }
export type ApiErrorBody = {
  error: { code: string; message: string };
};

export class ApiError extends Error {
  status: number;
  code: string;
  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

// apiFetch は現在のリクエストのCookieを Go へ転送して fetch する。
export async function apiFetch(
  path: string,
  init?: RequestInit,
): Promise<Response> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();

  return fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(cookieHeader ? { cookie: cookieHeader } : {}),
      ...init?.headers,
    },
    // 認証付きデータは常に最新を取得（キャッシュしない）。
    cache: "no-store",
  });
}

// apiJSON はJSONを返すエンドポイント用。非2xxは ApiError を投げる。
export async function apiJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await apiFetch(path, init);
  if (!res.ok) {
    throw await toApiError(res);
  }
  return (await res.json()) as T;
}

export async function toApiError(res: Response): Promise<ApiError> {
  try {
    const body = (await res.json()) as ApiErrorBody;
    return new ApiError(res.status, body.error.code, body.error.message);
  } catch {
    return new ApiError(res.status, "unknown", `リクエストに失敗しました (${res.status})`);
  }
}
