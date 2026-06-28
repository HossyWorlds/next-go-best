"use client";

import { deleteTaskAction, updateTaskAction } from "@/app/tasks/actions";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Task } from "@/lib/types/task";

import { StatusBadge } from "./status-badge";
import { TaskDialog } from "./task-dialog";

export function TaskItem({ task }: { task: Task }) {
  return (
    <Card>
      <CardHeader className="flex-row items-start justify-between gap-2">
        <CardTitle className="break-all">{task.title}</CardTitle>
        <StatusBadge status={task.status} />
      </CardHeader>
      <CardContent className="space-y-3">
        {task.description && (
          <p className="whitespace-pre-wrap text-sm text-muted-foreground">
            {task.description}
          </p>
        )}
        <div className="flex gap-2">
          <TaskDialog
            action={updateTaskAction}
            task={task}
            triggerLabel="編集"
            title="タスクを編集"
          />
          <form action={deleteTaskAction}>
            <input type="hidden" name="id" value={task.id} />
            <Button type="submit" variant="destructive" size="sm">
              削除
            </Button>
          </form>
        </div>
      </CardContent>
    </Card>
  );
}
