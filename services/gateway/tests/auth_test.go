package tests

import (
	"testing"
	"time"

	"gateway-service/internal/auth"
)

func TestIssueAndVerify(t *testing.T) {
	mgr := auth.New("test-secret", time.Hour)

	token, err := mgr.Issue(42)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := mgr.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected UserID=42, got %d", claims.UserID)
	}
}

func TestVerify_ExpiredToken(t *testing.T) {
	mgr := auth.New("test-secret", -time.Second)
	token, _ := mgr.Issue(1)

	_, err := mgr.Verify(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestVerify_WrongSecret(t *testing.T) {
	mgr1 := auth.New("secret-a", time.Hour)
	mgr2 := auth.New("secret-b", time.Hour)

	token, _ := mgr1.Issue(7)
	_, err := mgr2.Verify(token)
	if err == nil {
		t.Fatal("expected error when verifying with wrong secret")
	}
}

func TestVerify_InvalidToken(t *testing.T) {
	mgr := auth.New("test-secret", time.Hour)
	_, err := mgr.Verify("not.a.jwt")
	if err == nil {
		t.Fatal("expected error for garbage token")
	}
}

func TestVerify_EmptyToken(t *testing.T) {
	mgr := auth.New("test-secret", time.Hour)
	_, err := mgr.Verify("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestIssue_DifferentUsers(t *testing.T) {
	mgr := auth.New("test-secret", time.Hour)

	tok1, _ := mgr.Issue(1)
	tok2, _ := mgr.Issue(2)

	if tok1 == tok2 {
		t.Fatal("tokens for different users should differ (different iat)")
	}

	c1, _ := mgr.Verify(tok1)
	c2, _ := mgr.Verify(tok2)

	if c1.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", c1.UserID)
	}
	if c2.UserID != 2 {
		t.Errorf("expected UserID=2, got %d", c2.UserID)
	}
}

func TestIssue_IssuerField(t *testing.T) {
	mgr := auth.New("secret", time.Hour)
	token, _ := mgr.Issue(99)
	claims, _ := mgr.Verify(token)
	if claims.Issuer != "parking-gateway" {
		t.Errorf("expected issuer=parking-gateway, got %q", claims.Issuer)
	}
}
