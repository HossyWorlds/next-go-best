package auth

import (
	"net/http"
	"time"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
	"github.com/HossyWorlds/next-go-best/backend/internal/server"
)

// Handler は認証のHTTP境界。Cookieの発行/破棄もここで行う。
type Handler struct {
	svc          *Service
	cookieSecure bool
}

func NewHandler(svc *Service, cookieSecure bool) *Handler {
	return &Handler{svc: svc, cookieSecure: cookieSecure}
}

// RegisterRoutes は認証ルートを登録する。
// loginRateLimit はログインに適用するレート制限ミドルウェア、protect は認証必須ルート用。
func (h *Handler) RegisterRoutes(mux *http.ServeMux, loginRateLimit, protect func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/auth/register", http.HandlerFunc(h.Register))
	mux.Handle("POST /api/v1/auth/login", loginRateLimit(http.HandlerFunc(h.Login)))
	mux.Handle("POST /api/v1/auth/logout", http.HandlerFunc(h.Logout))
	mux.Handle("GET /api/v1/auth/me", protect(http.HandlerFunc(h.Me)))
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req credentials
	if err := server.DecodeJSON(r, &req); err != nil {
		server.WriteError(w, r, err)
		return
	}
	user, err := h.svc.Register(r.Context(), RegisterParams(req))
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	server.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req credentials
	if err := server.DecodeJSON(r, &req); err != nil {
		server.WriteError(w, r, err)
		return
	}
	res, err := h.svc.Login(r.Context(), LoginParams(req))
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	h.setSessionCookie(w, res.RawToken, res.ExpiresAt)
	server.WriteJSON(w, http.StatusOK, res.User)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(SessionCookieName); err == nil {
		if err := h.svc.Logout(r.Context(), c.Value); err != nil {
			server.WriteError(w, r, err)
			return
		}
	}
	h.clearSessionCookie(w)
	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		server.WriteError(w, r, apperr.Unauthorized("unauthorized", "認証が必要です"))
		return
	}
	server.WriteJSON(w, http.StatusOK, user)
}

// setSessionCookie は HttpOnly / SameSite=Lax / (本番)Secure のCookieを発行する。
func (h *Handler) setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}
