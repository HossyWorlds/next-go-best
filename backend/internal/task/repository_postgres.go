package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	dbgen "github.com/HossyWorlds/next-go-best/backend/internal/db/gen"
)

// PostgresRepository は sqlc 生成コードを使った Repository 実装。
type PostgresRepository struct {
	q *dbgen.Queries
}

// NewPostgresRepository は pgxpool.Pool などの DBTX を受け取る。
func NewPostgresRepository(db dbgen.DBTX) *PostgresRepository {
	return &PostgresRepository{q: dbgen.New(db)}
}

func (r *PostgresRepository) Create(ctx context.Context, in CreateInput) (Task, error) {
	row, err := r.q.CreateTask(ctx, dbgen.CreateTaskParams{
		OwnerID:     in.OwnerID,
		Title:       in.Title,
		Description: in.Description,
		Status:      string(in.Status),
	})
	if err != nil {
		return Task{}, err
	}
	return toDomain(row), nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (Task, error) {
	row, err := r.q.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}
	return toDomain(row), nil
}

func (r *PostgresRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int32) ([]Task, error) {
	rows, err := r.q.ListTasksByOwner(ctx, dbgen.ListTasksByOwnerParams{
		OwnerID: ownerID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]Task, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDomain(row))
	}
	return out, nil
}

func (r *PostgresRepository) CountByOwner(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	return r.q.CountTasksByOwner(ctx, ownerID)
}

func (r *PostgresRepository) Update(ctx context.Context, in UpdateInput) (Task, error) {
	row, err := r.q.UpdateTask(ctx, dbgen.UpdateTaskParams{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Status:      string(in.Status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}
	return toDomain(row), nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteTask(ctx, id)
}

func toDomain(row dbgen.Task) Task {
	return Task{
		ID:          row.ID,
		OwnerID:     row.OwnerID,
		Title:       row.Title,
		Description: row.Description,
		Status:      Status(row.Status),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
