package kit

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("admin123")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword("admin123", hash) {
		t.Fatal("expected password match")
	}
	if CheckPassword("wrong", hash) {
		t.Fatal("expected password mismatch")
	}
}

func TestSeedBcryptHash(t *testing.T) {
	// 与 00007 迁移种子保持一致
	const seed = "$2a$10$qiHGGPjc74psJ7JZTuyzX.0C6B7LYdD5SFoBDvdHhRXNIeLN/rp16"
	if !CheckPassword("admin123", seed) {
		t.Fatal("seed bcrypt hash should verify admin123")
	}
}
