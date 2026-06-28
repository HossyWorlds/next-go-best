package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
)

type ctxKey int

const userKey ctxKey = iota

// WithUser は認証済みユーザーをコンテキストに格納する（Middleware が使う）。
func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserFromContext は認証済みユーザーを取り出す。
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

// RequireUserID は認証済みユーザーIDを返す。未認証なら Unauthorized エラー。
// 保護ルート配下のハンドラから利用する。
func RequireUserID(ctx context.Context) (uuid.UUID, error) {
	u, ok := UserFromContext(ctx)
	if !ok {
		return uuid.Nil, apperr.Unauthorized("unauthorized", "認証が必要です")
	}
	return u.ID, nil
}
