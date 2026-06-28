import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

import type { TaskFormState } from "@/app/tasks/actions";
import type { Task } from "@/lib/types/task";

import { TaskDialog } from "./task-dialog";

const noopAction = async (): Promise<TaskFormState> => ({});

const sampleTask: Task = {
  id: "11111111-1111-1111-1111-111111111111",
  ownerId: "22222222-2222-2222-2222-222222222222",
  title: "既存タスク",
  description: "説明文",
  status: "doing",
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
};

describe("TaskDialog", () => {
  it("トリガーボタンのラベルを表示する", () => {
    render(
      <TaskDialog
        action={noopAction}
        triggerLabel="新規作成"
        title="タスクを作成"
      />,
    );
    expect(
      screen.getByRole("button", { name: "新規作成" }),
    ).toBeInTheDocument();
  });

  it("クリックでフォーム項目が開く", async () => {
    const user = userEvent.setup();
    render(
      <TaskDialog
        action={noopAction}
        triggerLabel="新規作成"
        title="タスクを作成"
      />,
    );
    await user.click(screen.getByRole("button", { name: "新規作成" }));
    expect(screen.getByLabelText("タイトル")).toBeInTheDocument();
    expect(screen.getByLabelText("ステータス")).toBeInTheDocument();
  });

  it("編集モードでは既存値が初期表示される", async () => {
    const user = userEvent.setup();
    render(
      <TaskDialog
        action={noopAction}
        triggerLabel="編集"
        title="タスクを編集"
        task={sampleTask}
      />,
    );
    await user.click(screen.getByRole("button", { name: "編集" }));
    expect(screen.getByLabelText("タイトル")).toHaveValue("既存タスク");
    expect(screen.getByLabelText("ステータス")).toHaveValue("doing");
  });
});
