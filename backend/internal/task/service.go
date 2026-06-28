package task

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
	"github.com/HossyWorlds/next-go-best/backend/internal/validate"
)

// ページネーションの既定値・上限。全件返却を避けるためのガード。
const (
	DefaultLimit int32 = 20
	MaxLimit     int32 = 100
)

// Service はタスクのビジネスロジック。永続化は Repository に委譲する。
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create はバリデーションのうえタスクを作成する。status 未指定は todo。
func (s *Service) Create(ctx context.Context, ownerID uuid.UUID, p CreateParams) (Task, error) {
	if p.Status == "" {
		p.Status = StatusTodo
	}
	if err := validate.Struct(p); err != nil {
		return Task{}, err
	}

	t, err := s.repo.Create(ctx, CreateInput{
		OwnerID:     ownerID,
		Title:       p.Title,
		Description: p.Description,
		Status:      p.Status,
	})
	if err != nil {
		return Task{}, apperr.Internal(err)
	}
	return t, nil
}

// Get は所有者本人のタスクのみ返す。他人のタスクは 403、無ければ 404。
func (s *Service) Get(ctx context.Context, userID, id uuid.UUID) (Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, mapGetError(err)
	}
	if t.OwnerID != userID {
		// 注: 情報漏えいを嫌うなら404に寄せる選択もある。ここでは権限の明示を優先し403。
		return Task{}, apperr.Forbidden("task_forbidden", "このタスクへのアクセス権がありません")
	}
	return t, nil
}

// List は所有者のタスクをページネーションして返す。
func (s *Service) List(ctx context.Context, ownerID uuid.UUID, limit, offset int32) (ListResult, error) {
	limit = normalizeLimit(limit)
	if offset < 0 {
		offset = 0
	}

	tasks, err := s.repo.ListByOwner(ctx, ownerID, limit, offset)
	if err != nil {
		return ListResult{}, apperr.Internal(err)
	}
	total, err := s.repo.CountByOwner(ctx, ownerID)
	if err != nil {
		return ListResult{}, apperr.Internal(err)
	}
	return ListResult{Tasks: tasks, Total: total}, nil
}

// Update は所有者本人のタスクのみ更新する。
func (s *Service) Update(ctx context.Context, userID, id uuid.UUID, p UpdateParams) (Task, error) {
	if err := validate.Struct(p); err != nil {
		return Task{}, err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, mapGetError(err)
	}
	if existing.OwnerID != userID {
		return Task{}, apperr.Forbidden("task_forbidden", "このタスクへのアクセス権がありません")
	}

	t, err := s.repo.Update(ctx, UpdateInput{
		ID:          id,
		Title:       p.Title,
		Description: p.Description,
		Status:      p.Status,
	})
	if err != nil {
		return Task{}, apperr.Internal(err)
	}
	return t, nil
}

// Delete は所有者本人のタスクのみ削除する。
func (s *Service) Delete(ctx context.Context, userID, id uuid.UUID) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return mapGetError(err)
	}
	if existing.OwnerID != userID {
		return apperr.Forbidden("task_forbidden", "このタスクへのアクセス権がありません")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func mapGetError(err error) error {
	if errors.Is(err, ErrNotFound) {
		return apperr.NotFound("task_not_found", "タスクが見つかりません")
	}
	return apperr.Internal(err)
}

func normalizeLimit(limit int32) int32 {
	switch {
	case limit <= 0:
		return DefaultLimit
	case limit > MaxLimit:
		return MaxLimit
	default:
		return limit
	}
}
