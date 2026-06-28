// Go API (internal/auth/model.go) の User DTO に対応するミラー型。
// password_hash は API レスポンスに含まれない（json:"-"）。
export type User = {
  id: string;
  email: string;
  createdAt: string;
  updatedAt: string;
};
