"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { logout } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { createTask, deleteTask, updateTask } from "@/lib/api/tasks";
import { TASK_STATUSES, type TaskStatus } from "@/lib/types/task";

export type TaskFormState = { error?: string; ok?: boolean };

export async function createTaskAction(
  _prev: TaskFormState,
  formData: FormData,
): Promise<TaskFormState> {
  try {
    await createTask({
      title: String(formData.get("title") ?? ""),
      description: String(formData.get("description") ?? ""),
      status: parseStatus(formData.get("status")),
    });
  } catch (e) {
    return { error: errorMessage(e, "作成に失敗しました") };
  }
  revalidatePath("/tasks");
  return { ok: true };
}

export async function updateTaskAction(
  _prev: TaskFormState,
  formData: FormData,
): Promise<TaskFormState> {
  const id = String(formData.get("id") ?? "");
  try {
    await updateTask(id, {
      title: String(formData.get("title") ?? ""),
      description: String(formData.get("description") ?? ""),
      status: parseStatus(formData.get("status")),
    });
  } catch (e) {
    return { error: errorMessage(e, "更新に失敗しました") };
  }
  revalidatePath("/tasks");
  return { ok: true };
}

// deleteTaskAction は <form action={...}> から直接呼ぶ（FormDataのみ）。
export async function deleteTaskAction(formData: FormData): Promise<void> {
  const id = String(formData.get("id") ?? "");
  await deleteTask(id);
  revalidatePath("/tasks");
}

export async function logoutAction(): Promise<void> {
  await logout();
  redirect("/login");
}

function parseStatus(value: FormDataEntryValue | null): TaskStatus {
  const v = String(value ?? "todo");
  return (TASK_STATUSES as readonly string[]).includes(v)
    ? (v as TaskStatus)
    : "todo";
}

function errorMessage(e: unknown, fallback: string): string {
  return e instanceof ApiError ? e.message : fallback;
}
