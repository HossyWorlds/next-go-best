import { expect, test } from "@playwright/test";

// コアフロー e2e: 登録 → ログイン → 作成 → 更新 → 削除 → ログアウト。
// 前提: Postgres + Go API(:8080) が起動していること（README参照）。

function uniqueEmail() {
  return `e2e-${Date.now()}-${Math.floor(Math.random() * 1e6)}@example.com`;
}

test("登録からCRUD・ログアウトまで通す", async ({ page }) => {
  const email = uniqueEmail();
  const password = "supersecret123";

  // 登録 → 自動ログイン → /tasks へ
  await page.goto("/register");
  await page.getByLabel("メールアドレス").fill(email);
  await page.getByLabel("パスワード").fill(password);
  await page.getByRole("button", { name: "登録する" }).click();

  await expect(page).toHaveURL(/\/tasks/);
  await expect(page.getByRole("heading", { name: "タスク" })).toBeVisible();
  await expect(page.getByText(email)).toBeVisible();

  // 作成
  await page.getByRole("button", { name: "新規作成" }).click();
  await page.getByLabel("タイトル").fill("E2Eタスク");
  await page.getByLabel("説明").fill("Playwrightから作成");
  await page.getByRole("button", { name: "保存" }).click();
  await expect(page.getByText("E2Eタスク")).toBeVisible();

  // 編集（ステータスを完了に）
  await page.getByRole("button", { name: "編集" }).click();
  await page.getByLabel("ステータス").selectOption("done");
  await page.getByRole("button", { name: "保存" }).click();
  // ステータスバッジが「完了」になっていること（select の option と区別する）。
  await expect(page.locator('[data-slot="badge"]', { hasText: "完了" })).toBeVisible();

  // 削除
  await page.getByRole("button", { name: "削除" }).click();
  await expect(page.getByText("E2Eタスク")).toHaveCount(0);

  // ログアウト → /login
  await page.getByRole("button", { name: "ログアウト" }).click();
  await expect(page).toHaveURL(/\/login/);
});

test("未ログインで /tasks は /login にリダイレクトされる", async ({ page }) => {
  await page.goto("/tasks");
  await expect(page).toHaveURL(/\/login/);
});
