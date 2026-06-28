package task

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
	"github.com/HossyWorlds/next-go-best/backend/internal/auth"
	"github.com/HossyWorlds/next-go-best/backend/internal/server"
)

// Handler はタスクのHTTP境界。デコード・認証ユーザー取得・サービス呼び出しに徹する。
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes は保護対象のタスクルートを mux に登録する。
// protect は認証ミドルウェア（auth.Middleware）を想定。
func (h *Handler) RegisterRoutes(mux *http.ServeMux, protect func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/tasks", protect(http.HandlerFunc(h.List)))
	mux.Handle("POST /api/v1/tasks", protect(http.HandlerFunc(h.Create)))
	mux.Handle("GET /api/v1/tasks/{id}", protect(http.HandlerFunc(h.Get)))
	mux.Handle("PUT /api/v1/tasks/{id}", protect(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /api/v1/tasks/{id}", protect(http.HandlerFunc(h.Delete)))
}

type writeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.RequireUserID(r.Context())
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	var req writeRequest
	if err := server.DecodeJSON(r, &req); err != nil {
		server.WriteError(w, r, err)
		return
	}

	created, err := h.svc.Create(r.Context(), userID, CreateParams(req))
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	server.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.RequireUserID(r.Context())
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	limit := parseInt32(r.URL.Query().Get("limit"), DefaultLimit)
	offset := parseInt32(r.URL.Query().Get("offset"), 0)

	res, err := h.svc.List(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	server.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, id, err := h.userAndID(r)
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	t, err := h.svc.Get(r.Context(), userID, id)
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	server.WriteJSON(w, http.StatusOK, t)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, id, err := h.userAndID(r)
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	var req writeRequest
	if err := server.DecodeJSON(r, &req); err != nil {
		server.WriteError(w, r, err)
		return
	}
	t, err := h.svc.Update(r.Context(), userID, id, UpdateParams(req))
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	server.WriteJSON(w, http.StatusOK, t)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, id, err := h.userAndID(r)
	if err != nil {
		server.WriteError(w, r, err)
		return
	}
	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		server.WriteError(w, r, err)
		return
	}
	server.WriteJSON(w, http.StatusNoContent, nil)
}

// userAndID は認証ユーザーIDとパスの {id} を取り出す共通処理。
func (h *Handler) userAndID(r *http.Request) (uuid.UUID, uuid.UUID, error) {
	userID, err := auth.RequireUserID(r.Context())
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return uuid.Nil, uuid.Nil, apperr.Validation("invalid_id", "IDの形式が不正です")
	}
	return userID, id, nil
}

func parseInt32(s string, fallback int32) int32 {
	if s == "" {
		return fallback
	}
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return fallback
	}
	return int32(n)
}
