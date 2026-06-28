"use client";

import { useState, useTransition } from "react";

import type { TaskFormState } from "@/app/tasks/actions";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { TASK_STATUSES, type Task } from "@/lib/types/task";

import { STATUS_LABEL } from "./status-badge";

type Props = {
  action: (prev: TaskFormState, fd: FormData) => Promise<TaskFormState>;
  triggerLabel: string;
  title: string;
  task?: Task; // 指定時は編集モード
};

// 作成・編集兼用のタスクフォーム（shadcn Dialog 内）。
// 成功（state.ok）でダイアログを閉じる。
export function TaskDialog({ action, triggerLabel, title, task }: Props) {
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string>();
  const [pending, startTransition] = useTransition();

  // サブミットはトランジション内で実行し、成功時にダイアログを閉じる。
  // （effect 内 setState を避けるため、結果はイベント内で処理する）
  function handleSubmit(formData: FormData) {
    startTransition(async () => {
      const result = await action({}, formData);
      if (result.error) {
        setError(result.error);
      } else {
        setError(undefined);
        setOpen(false);
      }
    });
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger
        render={
          <Button variant={task ? "outline" : "default"} size="sm">
            {triggerLabel}
          </Button>
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <form action={handleSubmit} className="space-y-4">
          {task && <input type="hidden" name="id" value={task.id} />}
          {error && (
            <p role="alert" className="text-sm text-destructive">
              {error}
            </p>
          )}
          <div className="space-y-2">
            <Label htmlFor="title">タイトル</Label>
            <Input
              id="title"
              name="title"
              defaultValue={task?.title}
              required
              maxLength={255}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">説明</Label>
            <Textarea
              id="description"
              name="description"
              defaultValue={task?.description}
              maxLength={2000}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="status">ステータス</Label>
            <select
              id="status"
              name="status"
              defaultValue={task?.status ?? "todo"}
              className="border-input flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
            >
              {TASK_STATUSES.map((s) => (
                <option key={s} value={s}>
                  {STATUS_LABEL[s]}
                </option>
              ))}
            </select>
          </div>
          <DialogFooter>
            <Button type="submit" disabled={pending}>
              {pending ? "保存中..." : "保存"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
