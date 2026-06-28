import type { Task } from "@/lib/types/task";

import { TaskItem } from "./task-item";

export function TaskList({ tasks }: { tasks: Task[] }) {
  if (tasks.length === 0) {
    return (
      <p className="text-sm text-muted-foreground" data-testid="empty-state">
        タスクはまだありません。「新規作成」から追加しましょう。
      </p>
    );
  }
  return (
    <div className="grid gap-3 sm:grid-cols-2">
      {tasks.map((task) => (
        <TaskItem key={task.id} task={task} />
      ))}
    </div>
  );
}
