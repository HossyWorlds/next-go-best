import { Badge } from "@/components/ui/badge";
import type { TaskStatus } from "@/lib/types/task";

const STATUS_LABEL: Record<TaskStatus, string> = {
  todo: "未着手",
  doing: "進行中",
  done: "完了",
};

const STATUS_VARIANT: Record<
  TaskStatus,
  "default" | "secondary" | "outline"
> = {
  todo: "outline",
  doing: "secondary",
  done: "default",
};

export function StatusBadge({ status }: { status: TaskStatus }) {
  return <Badge variant={STATUS_VARIANT[status]}>{STATUS_LABEL[status]}</Badge>;
}

export { STATUS_LABEL };
