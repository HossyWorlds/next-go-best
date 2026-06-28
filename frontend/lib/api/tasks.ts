import { apiFetch, apiJSON, toApiError } from "./client";
import type { Task, TaskList, TaskInput } from "@/lib/types/task";

// タスクAPIのラッパ。すべてサーバ側で実行され、Cookieを Go へ中継する。

export function listTasks(limit = 20, offset = 0): Promise<TaskList> {
  return apiJSON<TaskList>(`/api/v1/tasks?limit=${limit}&offset=${offset}`);
}

export function getTask(id: string): Promise<Task> {
  return apiJSON<Task>(`/api/v1/tasks/${id}`);
}

export function createTask(input: TaskInput): Promise<Task> {
  return apiJSON<Task>("/api/v1/tasks", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateTask(id: string, input: TaskInput): Promise<Task> {
  return apiJSON<Task>(`/api/v1/tasks/${id}`, {
    method: "PUT",
    body: JSON.stringify(input),
  });
}

export async function deleteTask(id: string): Promise<void> {
  const res = await apiFetch(`/api/v1/tasks/${id}`, { method: "DELETE" });
  if (!res.ok) {
    throw await toApiError(res);
  }
}
