-- tasks: ユーザー所有のCRUD対象。status はアプリ層の値域とDB制約の両方で守る。
CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'doing', 'done')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tasks_owner_id ON tasks (owner_id);
-- 一覧の既定ソート（新しい順）を効かせるための複合インデックス。
CREATE INDEX idx_tasks_owner_created ON tasks (owner_id, created_at DESC);
