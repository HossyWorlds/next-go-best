package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailTaken      = errors.New("email already taken")
	ErrSessionNotFound = errors.New("session not found")
)

// Repository は認証関連の永続化を抽象化する。
type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)

	CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (Session, error)
	// GetSessionByTokenHash は有効なセッションと所有ユーザーを返す。無効・期限切れは ErrSessionNotFound。
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (Session, User, error)
	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
}
