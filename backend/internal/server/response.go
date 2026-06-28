package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
)

// errorEnvelope は全エラーレスポンスの共通フォーマット。
//
//	{ "error": { "code": "task_not_found", "message": "..." } }
type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSON は成功レスポンスをJSONで書き出す。
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// ヘッダ送出済みのため、ここではログのみ。
		slog.Error("failed to encode json response", slog.Any("error", err))
	}
}

// WriteError は error を適切なHTTPステータス＋エラー封筒に変換して返す。
// *apperr.Error の Kind を見てマッピングし、未分類は 500 にフォールバックする。
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		appErr = apperr.Internal(err)
	}

	status := statusFromKind(appErr.Kind)

	// 5xx は原因をログに残す（レスポンス本文には出さない）。
	if status >= http.StatusInternalServerError {
		slog.LogAttrs(r.Context(), slog.LevelError, "request_error",
			slog.String("request_id", RequestIDFromContext(r.Context())),
			slog.String("code", appErr.Code),
			slog.Any("error", appErr.Err),
		)
	}

	writeErrorEnvelope(w, status, appErr.Code, appErr.Message)
}

func writeErrorEnvelope(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, errorEnvelope{Error: errorBody{Code: code, Message: message}})
}

func statusFromKind(kind apperr.Kind) int {
	switch kind {
	case apperr.KindValidation:
		return http.StatusBadRequest
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// DecodeJSON はリクエストボディをデコードする。失敗時は検証エラーを返す。
func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return apperr.Validation("invalid_json", "リクエストボディの形式が不正です")
	}
	return nil
}
