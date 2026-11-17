package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"
	"go.uber.org/mock/gomock"

	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/mocks"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	"github.com/go-eagle/eagle-layout/internal/types"
)

func TestUserServiceRegister(t *testing.T) {
	errRepo := errors.New("repo failed")
	errQueue := errors.New("queue unavailable")
	longPassword := strings.Repeat("a", 73)

	tests := []struct {
		name     string
		input    types.RegisterInput
		setup    func(t *testing.T, svc *userService, repo *mocks.MockUserRepo)
		wantErr  error
		wantID   int64
		wantName string
	}{
		{
			name:  "success",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, nil)
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, data model.UserInfoModel) (int64, error) {
						if data.Username != testUsername || data.Email != testEmail {
							t.Errorf("create user = %+v", data)
						}
						if !auth.ComparePasswords(data.Password, testPassword) {
							t.Error("password was not hashed")
						}
						return int64(42), nil
					},
				)
				svc.enqueueWelcome = func(payload tasks.EmailWelcomePayload) error {
					if payload.UserID != 42 {
						t.Errorf("welcome user id = %d, want 42", payload.UserID)
					}
					return nil
				}
			},
			wantID:   42,
			wantName: testUsername,
		},
		{
			name:  "email already exists",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(&model.UserInfoModel{ID: 1}, nil)
			},
			wantErr: ecode.ErrUserIsExist,
		},
		{
			name:  "username already exists",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(&model.UserInfoModel{ID: 2}, nil)
			},
			wantErr: ecode.ErrUserIsExist,
		},
		{
			name:  "email lookup fails",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:  "username lookup fails",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, errRepo)
			},
			wantErr: ecode.ErrInternalError,
		},
		{
			name:  "password cannot be hashed",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: longPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, nil)
			},
			wantErr: errcode.ErrEncrypt,
		},
		{
			name:  "create user fails",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, nil)
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:  "welcome task fails",
			input: types.RegisterInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, svc *userService, repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, nil)
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(7), nil)
				svc.enqueueWelcome = func(tasks.EmailWelcomePayload) error { return errQueue }
			},
			wantErr: errQueue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newMockedService(t)
			tt.setup(t, svc, repo)

			got, err := svc.Register(context.Background(), tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Register() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Register() unexpected error: %v", err)
			}
			if got.ID != tt.wantID || got.Username != tt.wantName {
				t.Fatalf("Register() = %+v, want id %d name %s", got, tt.wantID, tt.wantName)
			}
		})
	}
}

func TestUserServiceLogin(t *testing.T) {
	errRepo := errors.New("repo failed")
	errCache := errors.New("redis down")
	errSign := errors.New("sign failed")
	hashed := hashPassword(t, testPassword)
	user := &model.UserInfoModel{ID: 7, Username: testUsername, Password: hashed}

	tests := []struct {
		name    string
		input   types.LoginInput
		setup   func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache)
		wantErr error
		wantID  int64
	}{
		{
			name:  "success by email",
			input: types.LoginInput{Email: testEmail, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(user, nil)
				tokens.EXPECT().
					SetUserTokenCache(gomock.Any(), user.ID, gomock.Any(), cache.UserTokenExpireTime).
					Return(nil)
			},
			wantID: user.ID,
		},
		{
			name:  "email miss falls back to username",
			input: types.LoginInput{Email: testEmail, Username: testUsername, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, nil)
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(user, nil)
				tokens.EXPECT().
					SetUserTokenCache(gomock.Any(), user.ID, gomock.Any(), cache.UserTokenExpireTime).
					Return(nil)
			},
			wantID: user.ID,
		},
		{
			name:  "success by username only",
			input: types.LoginInput{Username: testUsername, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(user, nil)
				tokens.EXPECT().
					SetUserTokenCache(gomock.Any(), user.ID, gomock.Any(), cache.UserTokenExpireTime).
					Return(nil)
			},
			wantID: user.ID,
		},
		{
			name:    "missing credentials",
			input:   types.LoginInput{Password: testPassword},
			setup:   func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {},
			wantErr: ecode.ErrPasswordIncorrect,
		},
		{
			name:  "wrong password",
			input: types.LoginInput{Email: testEmail, Password: "nope"},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(user, nil)
			},
			wantErr: ecode.ErrPasswordIncorrect,
		},
		{
			name:  "user id is empty",
			input: types.LoginInput{Email: testEmail, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(&model.UserInfoModel{}, nil)
			},
			wantErr: ecode.ErrPasswordIncorrect,
		},
		{
			name:  "email lookup fails",
			input: types.LoginInput{Email: testEmail, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:  "username lookup fails",
			input: types.LoginInput{Username: testUsername, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:  "token signing fails",
			input: types.LoginInput{Email: testEmail, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(user, nil)
				svc.signToken = func(context.Context, map[string]interface{}, string, int64) (string, error) {
					return "", errSign
				}
			},
			wantErr: ecode.ErrToken,
		},
		{
			name:  "token cache write fails",
			input: types.LoginInput{Email: testEmail, Password: testPassword},
			setup: func(svc *userService, repo *mocks.MockUserRepo, tokens *mocks.MockUserTokenCache) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(user, nil)
				tokens.EXPECT().
					SetUserTokenCache(gomock.Any(), user.ID, gomock.Any(), cache.UserTokenExpireTime).
					Return(errCache)
			},
			wantErr: ecode.ErrToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, tokens := newMockedService(t)
			tt.setup(svc, repo, tokens)

			got, err := svc.Login(context.Background(), tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Login() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Login() unexpected error: %v", err)
			}
			if got.ID != tt.wantID || got.AccessToken == "" {
				t.Fatalf("Login() = %+v, want id %d and a token", got, tt.wantID)
			}
		})
	}
}

func TestUserServiceLogout(t *testing.T) {
	errCache := errors.New("redis down")
	const accessToken = "token-1"

	tests := []struct {
		name    string
		input   types.LogoutInput
		setup   func(tokens *mocks.MockUserTokenCache)
		wantErr error
	}{
		{
			name:  "success",
			input: types.LogoutInput{ID: 7, AccessToken: accessToken},
			setup: func(tokens *mocks.MockUserTokenCache) {
				tokens.EXPECT().GetUserTokenCache(gomock.Any(), int64(7)).Return(accessToken, nil)
				tokens.EXPECT().DelUserTokenCache(gomock.Any(), int64(7)).Return(nil)
			},
		},
		{
			name:  "token does not match",
			input: types.LogoutInput{ID: 7, AccessToken: accessToken},
			setup: func(tokens *mocks.MockUserTokenCache) {
				tokens.EXPECT().GetUserTokenCache(gomock.Any(), int64(7)).Return("other", nil)
			},
			wantErr: ecode.ErrAccessDenied,
		},
		{
			name:  "token lookup fails",
			input: types.LogoutInput{ID: 7, AccessToken: accessToken},
			setup: func(tokens *mocks.MockUserTokenCache) {
				tokens.EXPECT().GetUserTokenCache(gomock.Any(), int64(7)).Return("", errCache)
			},
			wantErr: ecode.ErrToken,
		},
		{
			name:  "token delete fails",
			input: types.LogoutInput{ID: 7, AccessToken: accessToken},
			setup: func(tokens *mocks.MockUserTokenCache) {
				tokens.EXPECT().GetUserTokenCache(gomock.Any(), int64(7)).Return(accessToken, nil)
				tokens.EXPECT().DelUserTokenCache(gomock.Any(), int64(7)).Return(errCache)
			},
			wantErr: ecode.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, tokens := newMockedService(t)
			tt.setup(tokens)

			got, err := svc.Logout(context.Background(), tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Logout() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Logout() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("Logout() = nil, want empty output")
			}
		})
	}
}
