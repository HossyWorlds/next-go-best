package task

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
)

// fakeRepo は Repository のインメモリ実装。service の単体テスト用。
type fakeRepo struct {
	tasks      map[uuid.UUID]Task
	failCreate bool
}

func newFakeRepo() *fakeRepo { return &fakeRepo{tasks: map[uuid.UUID]Task{}} }

func (f *fakeRepo) Create(_ context.Context, in CreateInput) (Task, error) {
	if f.failCreate {
		return Task{}, errors.New("boom")
	}
	t := Task{
		ID:          uuid.New(),
		OwnerID:     in.OwnerID,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
	}
	f.tasks[t.ID] = t
	return t, nil
}

func (f *fakeRepo) GetByID(_ context.Context, id uuid.UUID) (Task, error) {
	t, ok := f.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return t, nil
}

func (f *fakeRepo) ListByOwner(_ context.Context, ownerID uuid.UUID, limit, offset int32) ([]Task, error) {
	var out []Task
	for _, t := range f.tasks {
		if t.OwnerID == ownerID {
			out = append(out, t)
		}
	}
	// offset/limit の簡易適用（順序はテストでは問わない）
	if int(offset) >= len(out) {
		return []Task{}, nil
	}
	out = out[offset:]
	if int(limit) < len(out) {
		out = out[:limit]
	}
	return out, nil
}

func (f *fakeRepo) CountByOwner(_ context.Context, ownerID uuid.UUID) (int64, error) {
	var n int64
	for _, t := range f.tasks {
		if t.OwnerID == ownerID {
			n++
		}
	}
	return n, nil
}

func (f *fakeRepo) Update(_ context.Context, in UpdateInput) (Task, error) {
	t, ok := f.tasks[in.ID]
	if !ok {
		return Task{}, ErrNotFound
	}
	t.Title, t.Description, t.Status = in.Title, in.Description, in.Status
	f.tasks[in.ID] = t
	return t, nil
}

func (f *fakeRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(f.tasks, id)
	return nil
}

// kindOf は err から apperr.Kind を取り出す（テスト用）。
func kindOf(t *testing.T, err error) apperr.Kind {
	t.Helper()
	var ae *apperr.Error
	if !errors.As(err, &ae) {
		t.Fatalf("expected *apperr.Error, got %T (%v)", err, err)
	}
	return ae.Kind
}

func TestService_Create(t *testing.T) {
	owner := uuid.New()

	t.Run("status未指定はtodoになる", func(t *testing.T) {
		svc := NewService(newFakeRepo())
		got, err := svc.Create(context.Background(), owner, CreateParams{Title: "買い物"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Status != StatusTodo {
			t.Errorf("status = %q, want %q", got.Status, StatusTodo)
		}
		if got.OwnerID != owner {
			t.Errorf("ownerID = %v, want %v", got.OwnerID, owner)
		}
	})

	t.Run("titleが空ならvalidationエラー", func(t *testing.T) {
		svc := NewService(newFakeRepo())
		_, err := svc.Create(context.Background(), owner, CreateParams{Title: ""})
		if got := kindOf(t, err); got != apperr.KindValidation {
			t.Errorf("kind = %v, want Validation", got)
		}
	})

	t.Run("不正なstatusはvalidationエラー", func(t *testing.T) {
		svc := NewService(newFakeRepo())
		_, err := svc.Create(context.Background(), owner, CreateParams{Title: "x", Status: "unknown"})
		if got := kindOf(t, err); got != apperr.KindValidation {
			t.Errorf("kind = %v, want Validation", got)
		}
	})
}

func TestService_Get_Ownership(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	owner := uuid.New()
	other := uuid.New()

	created, err := svc.Create(context.Background(), owner, CreateParams{Title: "秘密のタスク"})
	if err != nil {
		t.Fatalf("setup create: %v", err)
	}

	t.Run("所有者は取得できる", func(t *testing.T) {
		got, err := svc.Get(context.Background(), owner, created.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != created.ID {
			t.Errorf("id mismatch")
		}
	})

	t.Run("他人は403", func(t *testing.T) {
		_, err := svc.Get(context.Background(), other, created.ID)
		if got := kindOf(t, err); got != apperr.KindForbidden {
			t.Errorf("kind = %v, want Forbidden", got)
		}
	})

	t.Run("存在しないIDは404", func(t *testing.T) {
		_, err := svc.Get(context.Background(), owner, uuid.New())
		if got := kindOf(t, err); got != apperr.KindNotFound {
			t.Errorf("kind = %v, want NotFound", got)
		}
	})
}

func TestService_Update_Delete_Ownership(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	owner := uuid.New()
	other := uuid.New()

	created, _ := svc.Create(context.Background(), owner, CreateParams{Title: "元タイトル"})

	t.Run("他人による更新は403", func(t *testing.T) {
		_, err := svc.Update(context.Background(), other, created.ID, UpdateParams{Title: "改ざん", Status: StatusDone})
		if got := kindOf(t, err); got != apperr.KindForbidden {
			t.Errorf("kind = %v, want Forbidden", got)
		}
	})

	t.Run("所有者は更新できる", func(t *testing.T) {
		got, err := svc.Update(context.Background(), owner, created.ID, UpdateParams{Title: "新タイトル", Status: StatusDoing})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Title != "新タイトル" || got.Status != StatusDoing {
			t.Errorf("update not applied: %+v", got)
		}
	})

	t.Run("他人による削除は403", func(t *testing.T) {
		err := svc.Delete(context.Background(), other, created.ID)
		if got := kindOf(t, err); got != apperr.KindForbidden {
			t.Errorf("kind = %v, want Forbidden", got)
		}
	})

	t.Run("所有者は削除できる", func(t *testing.T) {
		if err := svc.Delete(context.Background(), owner, created.ID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := svc.Get(context.Background(), owner, created.ID); kindOf(t, err) != apperr.KindNotFound {
			t.Errorf("expected NotFound after delete")
		}
	})
}

func TestService_List_Pagination(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	owner := uuid.New()
	for range 3 {
		if _, err := svc.Create(context.Background(), owner, CreateParams{Title: "t"}); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	res, err := svc.List(context.Background(), owner, 0, 0) // limit=0 → 既定値
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 {
		t.Errorf("total = %d, want 3", res.Total)
	}
	if len(res.Tasks) != 3 {
		t.Errorf("len = %d, want 3", len(res.Tasks))
	}
}
