package app_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

// createTask は1件作成し、生成されたタスクIDを返す。
func createTask(t *testing.T, c *http.Client, title string) map[string]any {
	t.Helper()
	status, body := doJSON(t, c, http.MethodPost, "/api/v1/tasks",
		map[string]any{"title": title})
	if status != http.StatusCreated {
		t.Fatalf("create task status=%d body=%s", status, body)
	}
	var task map[string]any
	if err := json.Unmarshal(body, &task); err != nil {
		t.Fatalf("unmarshal task: %v", err)
	}
	return task
}

func TestTask_RequiresAuth(t *testing.T) {
	c := newClient(t) // 未ログイン
	if status, _ := doJSON(t, c, http.MethodGet, "/api/v1/tasks", nil); status != http.StatusUnauthorized {
		t.Errorf("unauthenticated list status=%d want 401", status)
	}
	if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/tasks",
		map[string]any{"title": "x"}); status != http.StatusUnauthorized {
		t.Errorf("unauthenticated create status=%d want 401", status)
	}
}

func TestTask_CRUD(t *testing.T) {
	c := newClient(t)
	registerAndLogin(t, c, uniqueEmail("task-crud"), "supersecret123")

	// create（status未指定 → todo）
	created := createTask(t, c, "牛乳を買う")
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("created task has no id")
	}
	if created["status"] != "todo" {
		t.Errorf("default status=%v want todo", created["status"])
	}

	// get
	status, body := doJSON(t, c, http.MethodGet, "/api/v1/tasks/"+id, nil)
	if status != http.StatusOK {
		t.Fatalf("get status=%d body=%s", status, body)
	}

	// list（作成分が含まれ、totalが1以上）
	status, body = doJSON(t, c, http.MethodGet, "/api/v1/tasks", nil)
	if status != http.StatusOK {
		t.Fatalf("list status=%d body=%s", status, body)
	}
	var list struct {
		Tasks []map[string]any `json:"tasks"`
		Total int64            `json:"total"`
	}
	_ = json.Unmarshal(body, &list)
	if list.Total < 1 || len(list.Tasks) < 1 {
		t.Errorf("list total=%d len=%d want >=1", list.Total, len(list.Tasks))
	}

	// update
	status, body = doJSON(t, c, http.MethodPut, "/api/v1/tasks/"+id,
		map[string]any{"title": "牛乳とパンを買う", "description": "近所のスーパー", "status": "doing"})
	if status != http.StatusOK {
		t.Fatalf("update status=%d body=%s", status, body)
	}
	var updated map[string]any
	_ = json.Unmarshal(body, &updated)
	if updated["title"] != "牛乳とパンを買う" || updated["status"] != "doing" {
		t.Errorf("update not applied: %v", updated)
	}

	// delete
	if status, body := doJSON(t, c, http.MethodDelete, "/api/v1/tasks/"+id, nil); status != http.StatusNoContent {
		t.Fatalf("delete status=%d body=%s", status, body)
	}
	// 削除後の get は 404
	if status, _ := doJSON(t, c, http.MethodGet, "/api/v1/tasks/"+id, nil); status != http.StatusNotFound {
		t.Errorf("get after delete status=%d want 404", status)
	}
}

func TestTask_OwnershipIsolation(t *testing.T) {
	// userA が作成したタスクに userB はアクセスできない（403）。
	userA := newClient(t)
	registerAndLogin(t, userA, uniqueEmail("owner-a"), "supersecret123")
	created := createTask(t, userA, "Aの秘密タスク")
	id := created["id"].(string)

	userB := newClient(t)
	registerAndLogin(t, userB, uniqueEmail("owner-b"), "supersecret123")

	t.Run("他人のtaskはget不可(403)", func(t *testing.T) {
		if status, _ := doJSON(t, userB, http.MethodGet, "/api/v1/tasks/"+id, nil); status != http.StatusForbidden {
			t.Errorf("status=%d want 403", status)
		}
	})
	t.Run("他人のtaskはupdate不可(403)", func(t *testing.T) {
		if status, _ := doJSON(t, userB, http.MethodPut, "/api/v1/tasks/"+id,
			map[string]any{"title": "改ざん", "status": "done"}); status != http.StatusForbidden {
			t.Errorf("status=%d want 403", status)
		}
	})
	t.Run("他人のtaskはdelete不可(403)", func(t *testing.T) {
		if status, _ := doJSON(t, userB, http.MethodDelete, "/api/v1/tasks/"+id, nil); status != http.StatusForbidden {
			t.Errorf("status=%d want 403", status)
		}
	})
	t.Run("userBの一覧にAのtaskは出ない", func(t *testing.T) {
		_, body := doJSON(t, userB, http.MethodGet, "/api/v1/tasks", nil)
		var list struct {
			Total int64 `json:"total"`
		}
		_ = json.Unmarshal(body, &list)
		if list.Total != 0 {
			t.Errorf("userB total=%d want 0", list.Total)
		}
	})
}

func TestTask_InvalidID(t *testing.T) {
	c := newClient(t)
	registerAndLogin(t, c, uniqueEmail("task-badid"), "supersecret123")
	if status, _ := doJSON(t, c, http.MethodGet, "/api/v1/tasks/not-a-uuid", nil); status != http.StatusBadRequest {
		t.Errorf("invalid id status=%d want 400", status)
	}
}

func TestTask_ValidationError(t *testing.T) {
	c := newClient(t)
	registerAndLogin(t, c, uniqueEmail("task-val"), "supersecret123")
	// title空 → 400
	if status, _ := doJSON(t, c, http.MethodPost, "/api/v1/tasks",
		map[string]any{"title": ""}); status != http.StatusBadRequest {
		t.Errorf("empty title status=%d want 400", status)
	}
}
