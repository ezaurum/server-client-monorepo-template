package templateapp

import (
	"context"
	"templateapp/auth"
	"testing"
)

func TestLogin(t *testing.T) {
	redis := NewMockRedis()
	authService := auth.NewAuthService(redis, "jwt_signed", "jwt_unsigned", []byte("test_secret"))

	id := 1
	login, err := authService.Login(context.Background(), uint(id))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if login == "" {
		t.Fatalf("Login returned empty token")
	}
}
