package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrNotFound はレコードが存在しないことを表す。repository が返し、service が解釈する。
var ErrNotFound = errors.New("task not found")

// CreateInput は永続化層への作成入力（所有者を含む）。
type CreateInput struct {
	OwnerID     uuid.UUID
	Title       string
	Description string
	Status      Status
}

// UpdateInput は永続化層への更新入力。
type UpdateInput struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      Status
}

// Repository はタスクの永続化を抽象化する。
// service はこのインターフェースにのみ依存し、テストでは fake に差し替える。
type Repository interface {
	Create(ctx context.Context, in CreateInput) (Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (Task, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int32) ([]Task, error)
	CountByOwner(ctx context.Context, ownerID uuid.UUID) (int64, error)
	Update(ctx context.Context, in UpdateInput) (Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
