package auth

import (
	"testing"
	"time"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := HashPassword("s3cret-pw")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "s3cret-pw" {
		t.Fatal("password was not hashed")
	}
	if !CheckPassword(hash, "s3cret-pw") {
		t.Error("correct password rejected")
	}
	if CheckPassword(hash, "wrong-pw") {
		t.Error("wrong password accepted")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	m := NewTokenManager("test-secret", time.Hour)

	tok, err := m.Generate(42)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	id, err := m.Verify(tok)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if id != 42 {
		t.Errorf("got user id %d, want 42", id)
	}
}

func TestTokenExpired(t *testing.T) {
	m := NewTokenManager("test-secret", -time.Minute) // already expired
	tok, err := m.Generate(1)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := m.Verify(tok); err == nil {
		t.Error("expected expired token to be rejected")
	}
}

func TestTokenWrongSecret(t *testing.T) {
	issuer := NewTokenManager("secret-a", time.Hour)
	tok, _ := issuer.Generate(1)

	verifier := NewTokenManager("secret-b", time.Hour)
	if _, err := verifier.Verify(tok); err == nil {
		t.Error("token signed with a different secret should be rejected")
	}
}
