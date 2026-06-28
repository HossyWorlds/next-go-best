package auth

import (
	"net/http"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
	"github.com/HossyWorlds/next-go-best/backend/internal/server"
)

// SessionCookieName はセッションCookieの名前。
const SessionCookieName = "session"

// Middleware はセッションCookieを検証し、ユーザーをコンテキストに注入する。
// 無効なら 401 を返してハンドラには到達させない（保護ルート用）。
func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookieName)
		if err != nil {
			server.WriteError(w, r, apperr.Unauthorized("unauthorized", "認証が必要です"))
			return
		}
		user, err := s.Authenticate(r.Context(), c.Value)
		if err != nil {
			server.WriteError(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
}
