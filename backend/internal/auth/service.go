package auth

import (
	"context"
	"errors"
	"time"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
	"github.com/HossyWorlds/next-go-best/backend/internal/validate"
)

// LoginResult はログイン成功時の戻り。RawToken は Cookie に載せる生トークン。
type LoginResult struct {
	User      User
	RawToken  string
	ExpiresAt time.Time
}

// Service は認証のビジネスロジック。
type Service struct {
	repo       Repository
	sessionTTL time.Duration
}

func NewService(repo Repository, sessionTTL time.Duration) *Service {
	return &Service{repo: repo, sessionTTL: sessionTTL}
}

// Register は新規ユーザーを作成する。email 重複は 409。
func (s *Service) Register(ctx context.Context, p RegisterParams) (User, error) {
	if err := validate.Struct(p); err != nil {
		return User{}, err
	}

	hash, err := HashPassword(p.Password)
	if err != nil {
		return User{}, apperr.Internal(err)
	}

	user, err := s.repo.CreateUser(ctx, p.Email, hash)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return User{}, apperr.Conflict("email_taken", "このメールアドレスは既に登録されています")
		}
		return User{}, apperr.Internal(err)
	}
	return user, nil
}

// Login は認証情報を検証し、セッションを発行する。
// 失敗時は user有無を問わず同一の Unauthorized を返す（ユーザー列挙対策）。
func (s *Service) Login(ctx context.Context, p LoginParams) (LoginResult, error) {
	if err := validate.Struct(p); err != nil {
		return LoginResult{}, err
	}

	user, err := s.repo.GetUserByEmail(ctx, p.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// タイミング差を抑えるためダミー検証を実行してから失敗を返す。
			_, _ = VerifyPassword(p.Password, dummyHash)
			return LoginResult{}, invalidCredentials()
		}
		return LoginResult{}, apperr.Internal(err)
	}

	ok, err := VerifyPassword(p.Password, user.PasswordHash)
	if err != nil {
		return LoginResult{}, apperr.Internal(err)
	}
	if !ok {
		return LoginResult{}, invalidCredentials()
	}

	rawToken, err := generateToken()
	if err != nil {
		return LoginResult{}, apperr.Internal(err)
	}
	expiresAt := time.Now().Add(s.sessionTTL)

	if _, err := s.repo.CreateSession(ctx, user.ID, hashToken(rawToken), expiresAt); err != nil {
		return LoginResult{}, apperr.Internal(err)
	}

	return LoginResult{User: user, RawToken: rawToken, ExpiresAt: expiresAt}, nil
}

// Logout はセッションを失効させる（冪等）。
func (s *Service) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}
	if err := s.repo.DeleteSessionByTokenHash(ctx, hashToken(rawToken)); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

// Authenticate は生トークンからユーザーを解決する。無効なら Unauthorized。
func (s *Service) Authenticate(ctx context.Context, rawToken string) (*User, error) {
	if rawToken == "" {
		return nil, apperr.Unauthorized("unauthorized", "認証が必要です")
	}
	_, user, err := s.repo.GetSessionByTokenHash(ctx, hashToken(rawToken))
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, apperr.Unauthorized("unauthorized", "セッションが無効です")
		}
		return nil, apperr.Internal(err)
	}
	return &user, nil
}

func invalidCredentials() error {
	return apperr.Unauthorized("invalid_credentials", "メールアドレスまたはパスワードが正しくありません")
}

// dummyHash はユーザー不在時のタイミング均一化用の固定 argon2id ハッシュ（"x"）。
const dummyHash = "$argon2id$v=19$m=65536,t=3,p=2$AAAAAAAAAAAAAAAAAAAAAA$RdescudvJCsgt3ub+b+dWRWJTmaaJObG"
