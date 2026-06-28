// Package migrations はマイグレーションSQLを埋め込み、バイナリ単体で適用可能にする。
// cmd/migrate と統合テスト（testcontainers）の双方から同じSQLを再利用する。
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
