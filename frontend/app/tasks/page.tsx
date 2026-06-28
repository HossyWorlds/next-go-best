import { redirect } from "next/navigation";

import { TaskDialog } from "@/components/tasks/task-dialog";
import { TaskList } from "@/components/tasks/task-list";
import { Button } from "@/components/ui/button";
import { getMe } from "@/lib/api/auth";
import { listTasks } from "@/lib/api/tasks";

import { createTaskAction, logoutAction } from "./actions";

// Server Component: 認証ユーザーのタスク一覧をサーバ側で取得して描画する。
export default async function TasksPage() {
  const me = await getMe();
  if (!me) {
    // Cookieはあるがセッション無効などのケース。
    redirect("/login");
  }

  const list = await listTasks();

  return (
    <main className="mx-auto max-w-3xl space-y-6 p-6">
      <header className="flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">タスク</h1>
          <p className="text-sm text-muted-foreground">{me.email}</p>
        </div>
        <div className="flex items-center gap-2">
          <TaskDialog
            action={createTaskAction}
            triggerLabel="新規作成"
            title="タスクを作成"
          />
          <form action={logoutAction}>
            <Button type="submit" variant="outline" size="sm">
              ログアウト
            </Button>
          </form>
        </div>
      </header>

      <p className="text-sm text-muted-foreground">全 {list.total} 件</p>
      <TaskList tasks={list.tasks} />
    </main>
  );
}
