import { defineConfig, devices } from "@playwright/test";

// e2e は実ブラウザでコアフロー（登録→ログイン→CRUD→ログアウト）を検証する。
// 前提: Postgres と Go API(:8080) が起動していること（README参照）。
// Next dev サーバ(:3000) は webServer で自動起動する。
export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: "list",
  use: {
    baseURL: "http://localhost:3000",
    trace: "on-first-retry",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer: {
    command: "pnpm dev",
    url: "http://localhost:3000",
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
});
