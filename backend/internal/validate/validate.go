// Package validate は go-playground/validator を薄くラップし、
// 検証失敗を apperr.Validation に変換して各レイヤから使えるようにする。
package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/HossyWorlds/next-go-best/backend/internal/apperr"
)

var v = validator.New()

// Struct は構造体のタグ検証を行い、失敗時は *apperr.Error(Validation) を返す。
func Struct(s any) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var verrs validator.ValidationErrors
	if !asValidationErrors(err, &verrs) {
		return apperr.Internal(err)
	}

	msgs := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		msgs = append(msgs, fieldMessage(fe))
	}
	return apperr.Validation("validation_error", strings.Join(msgs, "; "))
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	verrs, ok := err.(validator.ValidationErrors)
	if ok {
		*target = verrs
	}
	return ok
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s は必須です", fe.Field())
	case "email":
		return fmt.Sprintf("%s はメールアドレス形式である必要があります", fe.Field())
	case "min":
		return fmt.Sprintf("%s は %s 文字以上である必要があります", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s は %s 文字以下である必要があります", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s は次のいずれかである必要があります: %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s が不正です", fe.Field())
	}
}
