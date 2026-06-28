// Go API (internal/task/model.go) のレスポンスDTOに対応するミラー型。
// API契約が変わったら api-contract-sync skill で同期する。

export const TASK_STATUSES = ["todo", "doing", "done"] as const;
export type TaskStatus = (typeof TASK_STATUSES)[number];

export type Task = {
  id: string;
  ownerId: string;
  title: string;
  description: string;
  status: TaskStatus;
  createdAt: string;
  updatedAt: string;
};

// GET /api/v1/tasks のレスポンス（ListResult）。
export type TaskList = {
  tasks: Task[];
  total: number;
};

// 作成・更新の入力。
export type TaskInput = {
  title: string;
  description: string;
  status: TaskStatus;
};
