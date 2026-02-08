package utils

import (
	"os"
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(1, "tester", 2)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}

	if claims.UserID != 1 || claims.Username != "tester" || claims.RoleID != 2 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestValidateTokenMissingSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := ValidateToken("dummy")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	}
}
