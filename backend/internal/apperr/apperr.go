// Package apperr はドメイン層が返すアプリケーションエラーを定義する。
// HTTP やDBの詳細に依存せず、種別(Kind)とユーザー向けメッセージを持つ。
// HTTPステータスへの変換は server 層（response.go）が一手に引き受ける。
package apperr

import "fmt"

// Kind はエラーの分類。HTTPステータスへのマッピングに使う。
type Kind int

const (
	KindInternal     Kind = iota // 500
	KindValidation               // 400
	KindUnauthorized             // 401
	KindForbidden                // 403
	KindNotFound                 // 404
	KindConflict                 // 409
)

// Error はアプリケーションエラー。Code は機械可読な短い識別子。
type Error struct {
	Kind    Kind
	Code    string // 例: "task_not_found"
	Message string // ユーザー向けの説明（機密を含めない）
	Err     error  // 元エラー（ログ用。レスポンスには出さない）
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return e.Code
}

func (e *Error) Unwrap() error { return e.Err }

// wrap は元エラーを保持したまま新しい *Error を返す。
func newErr(kind Kind, code, msg string, cause error) *Error {
	return &Error{Kind: kind, Code: code, Message: msg, Err: cause}
}

func Internal(cause error) *Error {
	return newErr(KindInternal, "internal_error", "予期しないエラーが発生しました", cause)
}

func Validation(code, msg string) *Error {
	return newErr(KindValidation, code, msg, nil)
}

func Unauthorized(code, msg string) *Error {
	return newErr(KindUnauthorized, code, msg, nil)
}

func Forbidden(code, msg string) *Error {
	return newErr(KindForbidden, code, msg, nil)
}

func NotFound(code, msg string) *Error {
	return newErr(KindNotFound, code, msg, nil)
}

func Conflict(code, msg string) *Error {
	return newErr(KindConflict, code, msg, nil)
}
