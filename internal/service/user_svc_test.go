package service

import (
	"testing"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"go.uber.org/mock/gomock"

	"github.com/go-eagle/eagle-layout/internal/mocks"
	"github.com/go-eagle/eagle-layout/internal/tasks"
)

const (
	testJWTSecret = "unit-test-secret"
	testPassword  = "password123"
	testUsername  = "alice"
	testEmail     = "alice@example.com"
)

func newMockedService(t *testing.T) (*userService, *mocks.MockUserRepo, *mocks.MockUserTokenCache) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockUserRepo(ctrl)
	tokens := mocks.NewMockUserTokenCache(ctrl)
	svc := newUserService(repo, tokens, func(tasks.EmailWelcomePayload) error { return nil }, testJWTSecret)
	return svc, repo, tokens
}

func hashPassword(t *testing.T, plain string) string {
	t.Helper()
	hashed, err := auth.HashAndSalt(plain)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return hashed
}

func TestNewUserServiceUsesDefaultEnqueuer(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockUserRepo(ctrl)
	tokens := mocks.NewMockUserTokenCache(ctrl)

	svc := NewUserService(repo, tokens)
	impl, ok := svc.(*userService)
	if !ok || impl.enqueueWelcome == nil || impl.signToken == nil || impl.tokens == nil {
		t.Fatalf("NewUserService() = %#v, want a userService with default dependencies", svc)
	}
}

func TestTokenSecret(t *testing.T) {
	svc := &userService{jwtSecret: "override"}
	if got := svc.tokenSecret(); got != "override" {
		t.Fatalf("tokenSecret() = %q, want override", got)
	}

	svc.jwtSecret = ""
	previous := app.Conf
	t.Cleanup(func() { app.Conf = previous })

	app.Conf = nil
	if got := svc.tokenSecret(); got != "" {
		t.Fatalf("tokenSecret() with nil config = %q, want empty", got)
	}

	app.Conf = &app.Config{JwtSecret: "from-config"}
	if got := svc.tokenSecret(); got != "from-config" {
		t.Fatalf("tokenSecret() = %q, want from-config", got)
	}
}
