-- name: CreateTask :one
INSERT INTO tasks (owner_id, title, description, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasksByOwner :many
-- 一覧は必ず owner で絞り、LIMIT/OFFSET でページネーションする（全件返却しない）。
SELECT * FROM tasks
WHERE owner_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTasksByOwner :one
SELECT count(*) FROM tasks
WHERE owner_id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2, description = $3, status = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
