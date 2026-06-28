import { NextResponse, type NextRequest } from "next/server";

// Next.js 16 で middleware.ts は proxy.ts に改名された（関数名も proxy、nodejsランタイム固定）。
// ここではセッションCookieの有無だけで保護ルートを制御する。
// （トークンの正当性検証は Go API 側で行う。proxy は軽量な振り分けに徹する）
const SESSION_COOKIE = "session";

export function proxy(request: NextRequest) {
  const hasSession = request.cookies.has(SESSION_COOKIE);
  const { pathname } = request.nextUrl;

  // 未ログインで保護ルートにアクセス → /login へ
  if (!hasSession && pathname.startsWith("/tasks")) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  // ログイン済みで認証画面に来たら → /tasks へ
  if (hasSession && (pathname === "/login" || pathname === "/register")) {
    const url = request.nextUrl.clone();
    url.pathname = "/tasks";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/tasks/:path*", "/login", "/register"],
};
