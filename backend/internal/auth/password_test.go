package auth

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	const pw = "correct horse battery staple"

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == pw {
		t.Fatal("hash must not equal plaintext")
	}

	t.Run("正しいパスワードは一致", func(t *testing.T) {
		ok, err := VerifyPassword(pw, hash)
		if err != nil {
			t.Fatalf("VerifyPassword: %v", err)
		}
		if !ok {
			t.Error("expected match")
		}
	})

	t.Run("誤ったパスワードは不一致", func(t *testing.T) {
		ok, err := VerifyPassword("wrong", hash)
		if err != nil {
			t.Fatalf("VerifyPassword: %v", err)
		}
		if ok {
			t.Error("expected mismatch")
		}
	})

	t.Run("毎回ソルトが異なる", func(t *testing.T) {
		hash2, err := HashPassword(pw)
		if err != nil {
			t.Fatalf("HashPassword: %v", err)
		}
		if hash == hash2 {
			t.Error("hashes must differ due to random salt")
		}
	})

	t.Run("不正な形式はエラー", func(t *testing.T) {
		if _, err := VerifyPassword(pw, "not-a-valid-hash"); err == nil {
			t.Error("expected error for malformed hash")
		}
	})
}

func TestDummyHashIsValid(t *testing.T) {
	// ログインのタイミング均一化に使う dummyHash がパース可能であることを保証する。
	if _, err := VerifyPassword("anything", dummyHash); err != nil {
		t.Fatalf("dummyHash must be parseable: %v", err)
	}
}
