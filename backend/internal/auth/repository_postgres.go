package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	dbgen "github.com/HossyWorlds/next-go-best/backend/internal/db/gen"
)

// uniqueViolation は Postgres の一意制約違反コード。
const uniqueViolation = "23505"

// PostgresRepository は sqlc 生成コードを使った認証 Repository 実装。
type PostgresRepository struct {
	q *dbgen.Queries
}

func NewPostgresRepository(db dbgen.DBTX) *PostgresRepository {
	return &PostgresRepository{q: dbgen.New(db)}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	row, err := r.q.CreateUser(ctx, dbgen.CreateUserParams{Email: email, PasswordHash: passwordHash})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return User{}, ErrEmailTaken
		}
		return User{}, err
	}
	return userToDomain(row), nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return userToDomain(row), nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return userToDomain(row), nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (Session, error) {
	row, err := r.q.CreateSession(ctx, dbgen.CreateSessionParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return Session{}, err
	}
	return sessionToDomain(row), nil
}

func (r *PostgresRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (Session, User, error) {
	row, err := r.q.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Session{}, User{}, ErrSessionNotFound
		}
		return Session{}, User{}, err
	}
	return sessionToDomain(row.Session), userToDomain(row.User), nil
}

func (r *PostgresRepository) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	return r.q.DeleteSessionByTokenHash(ctx, tokenHash)
}

func userToDomain(u dbgen.User) User {
	return User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func sessionToDomain(s dbgen.Session) Session {
	return Session{
		ID:        s.ID,
		UserID:    s.UserID,
		TokenHash: s.TokenHash,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}
}
