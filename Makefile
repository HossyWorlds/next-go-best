# next-go-best 開発用ショートカット
# `.env` があれば読み込む（DATABASE_URL 等）
ifneq (,$(wildcard .env))
include .env
export
endif

.DEFAULT_GOAL := help

.PHONY: help up down logs migrate migrate-down sqlc dev-api test test-back lint-back \
        front-install dev-front test-front e2e build-front lint-front

help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

## ----- インフラ -----
up: ## Postgres を起動
	docker compose up -d db

down: ## Postgres を停止
	docker compose down

logs: ## Postgres のログを追う
	docker compose logs -f db

## ----- バックエンド -----
migrate: ## DBマイグレーションを適用（up）
	cd backend && go run ./cmd/migrate up

migrate-down: ## マイグレーションを1つ戻す（down）
	cd backend && go run ./cmd/migrate down

sqlc: ## sqlc でクエリからGoコードを生成
	cd backend && go tool sqlc generate

dev-api: ## APIサーバを起動
	cd backend && go run ./cmd/api

test: test-back test-front ## 全テスト

test-back: ## バックエンドのテスト（単体+統合 testcontainers）
	cd backend && go test ./...

lint-back: ## バックエンドの lint
	cd backend && go tool golangci-lint run

## ----- フロントエンド -----
front-install: ## 依存インストール
	cd frontend && pnpm install

dev-front: ## Next.js dev サーバ
	cd frontend && pnpm dev

test-front: ## Vitest コンポーネントテスト
	cd frontend && pnpm test

e2e: ## Playwright e2e
	cd frontend && pnpm exec playwright test

build-front: ## 本番ビルド
	cd frontend && pnpm build

lint-front: ## ESLint
	cd frontend && pnpm lint
