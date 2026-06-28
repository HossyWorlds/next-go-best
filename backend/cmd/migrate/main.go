// Command migrate は埋め込みSQLでDBマイグレーションを適用する。
//
//	go run ./cmd/migrate up    # 最新まで適用
//	go run ./cmd/migrate down  # 1ステップ戻す
package main

import (
	"log"
	"os"

	"github.com/HossyWorlds/next-go-best/backend/internal/db"
)

func main() {
	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	switch direction {
	case "up":
		if err := db.Migrate(databaseURL); err != nil {
			log.Fatalf("migrate up failed: %v", err)
		}
		log.Println("migrate up: done")
	case "down":
		if err := db.MigrateDown(databaseURL); err != nil {
			log.Fatalf("migrate down failed: %v", err)
		}
		log.Println("migrate down: done")
	default:
		log.Fatalf("unknown direction %q (use 'up' or 'down')", direction)
	}
}
