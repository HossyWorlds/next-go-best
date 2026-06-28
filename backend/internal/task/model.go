package task

import (
	"time"

	"github.com/google/uuid"
)

// Status はタスクの状態。DBの CHECK 制約と値域を一致させる。
type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

// Task はドメインモデル。
type Task struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"ownerId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateParams は新規作成の入力。検証タグは service が利用する。
type CreateParams struct {
	Title       string `validate:"required,max=255"`
	Description string `validate:"max=2000"`
	Status      Status `validate:"omitempty,oneof=todo doing done"`
}

// UpdateParams は更新の入力。
type UpdateParams struct {
	Title       string `validate:"required,max=255"`
	Description string `validate:"max=2000"`
	Status      Status `validate:"required,oneof=todo doing done"`
}

// ListResult は一覧とページネーション情報。
type ListResult struct {
	Tasks []Task `json:"tasks"`
	Total int64  `json:"total"`
}
