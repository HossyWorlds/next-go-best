package auth

import (
	"time"

	"github.com/google/uuid"
)

// User は認証ユーザー。password_hash はAPIレスポンスに含めない（json:"-"）。
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Session はサーバサイドセッション。生トークンは保持せずハッシュのみ。
type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// RegisterParams / LoginParams は入力DTO。検証タグは service が利用する。
type RegisterParams struct {
	Email    string `validate:"required,email,max=255"`
	Password string `validate:"required,min=8,max=72"` // bcrypt互換の上限に合わせ72
}

type LoginParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
