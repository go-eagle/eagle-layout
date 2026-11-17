package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"
	"go.uber.org/mock/gomock"

	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/mocks"
	"github.com/go-eagle/eagle-layout/internal/types"
)

func TestUserServiceCreateUser(t *testing.T) {
	errRepo := errors.New("repo failed")
	longPassword := strings.Repeat("a", 73)

	tests := []struct {
		name    string
		input   types.CreateUserInput
		setup   func(t *testing.T, repo *mocks.MockUserRepo)
		wantErr error
		wantID  int64
	}{
		{
			name:  "success",
			input: types.CreateUserInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, repo *mocks.MockUserRepo) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, data model.UserInfoModel) (int64, error) {
						if data.Username != testUsername || data.Email != testEmail {
							t.Errorf("create user = %+v", data)
						}
						if !auth.ComparePasswords(data.Password, testPassword) {
							t.Error("password was not hashed")
						}
						return int64(11), nil
					},
				)
			},
			wantID: 11,
		},
		{
			name:    "password cannot be hashed",
			input:   types.CreateUserInput{Username: testUsername, Email: testEmail, Password: longPassword},
			setup:   func(t *testing.T, repo *mocks.MockUserRepo) {},
			wantErr: errcode.ErrEncrypt,
		},
		{
			name:  "repository fails",
			input: types.CreateUserInput{Username: testUsername, Email: testEmail, Password: testPassword},
			setup: func(t *testing.T, repo *mocks.MockUserRepo) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newMockedService(t)
			tt.setup(t, repo)

			got, err := svc.CreateUser(context.Background(), tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("CreateUser() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("CreateUser() unexpected error: %v", err)
			}
			if got.ID != tt.wantID || got.Username != testUsername {
				t.Fatalf("CreateUser() = %+v, want id %d", got, tt.wantID)
			}
		})
	}
}

func TestUserServiceUpdateUser(t *testing.T) {
	errRepo := errors.New("repo failed")
	input := types.UpdateUserInput{
		UserId:   9,
		Nickname: "ally",
		Email:    testEmail,
		Avatar:   "a.png",
		Birthday: "2000-01-01",
		Bio:      "hi",
		Status:   1,
		Gender:   2,
	}

	tests := []struct {
		name    string
		input   types.UpdateUserInput
		setup   func(t *testing.T, repo *mocks.MockUserRepo)
		wantErr error
	}{
		{
			name:    "missing user id",
			input:   types.UpdateUserInput{Nickname: "ally"},
			setup:   func(t *testing.T, repo *mocks.MockUserRepo) {},
			wantErr: ecode.ErrInvalidArgument,
		},
		{
			name:  "success",
			input: input,
			setup: func(t *testing.T, repo *mocks.MockUserRepo) {
				repo.EXPECT().UpdateUser(gomock.Any(), input.UserId, gomock.Any()).DoAndReturn(
					func(_ context.Context, _ int64, data model.UserInfoModel) error {
						if data.Nickname != input.Nickname || data.Email != input.Email || data.Bio != input.Bio {
							t.Errorf("update data = %+v", data)
						}
						if data.UpdatedAt <= 0 {
							t.Error("updated_at was not set")
						}
						return nil
					},
				)
			},
		},
		{
			name:  "repository fails",
			input: input,
			setup: func(t *testing.T, repo *mocks.MockUserRepo) {
				repo.EXPECT().UpdateUser(gomock.Any(), input.UserId, gomock.Any()).Return(errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newMockedService(t)
			tt.setup(t, repo)

			got, err := svc.UpdateUser(context.Background(), tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdateUser() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateUser() unexpected error: %v", err)
			}
			if got.UserId != tt.input.UserId || got.Nickname != tt.input.Nickname || got.UpdatedAt <= 0 {
				t.Fatalf("UpdateUser() = %+v", got)
			}
		})
	}
}

func TestUserServiceUpdatePassword(t *testing.T) {
	errRepo := errors.New("repo failed")
	hashed := hashPassword(t, testPassword)
	longPassword := strings.Repeat("b", 73)
	stored := &model.UserInfoModel{ID: 3, Password: hashed}

	tests := []struct {
		name    string
		input   types.UpdatePasswordInput
		setup   func(repo *mocks.MockUserRepo)
		wantErr error
	}{
		{
			name:    "missing id",
			input:   types.UpdatePasswordInput{Password: testPassword, ConfirmPassword: testPassword, NewPassword: "next-password"},
			setup:   func(repo *mocks.MockUserRepo) {},
			wantErr: ecode.ErrInvalidArgument,
		},
		{
			name:    "empty password field",
			input:   types.UpdatePasswordInput{ID: 3, ConfirmPassword: testPassword, NewPassword: "next-password"},
			setup:   func(repo *mocks.MockUserRepo) {},
			wantErr: ecode.ErrInvalidArgument,
		},
		{
			name:    "confirmation does not match",
			input:   types.UpdatePasswordInput{ID: 3, Password: testPassword, ConfirmPassword: "other", NewPassword: "next-password"},
			setup:   func(repo *mocks.MockUserRepo) {},
			wantErr: ecode.ErrTwicePasswordNotMatch,
		},
		{
			name:  "lookup fails",
			input: types.UpdatePasswordInput{ID: 3, Password: testPassword, ConfirmPassword: testPassword, NewPassword: "next-password"},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(3)).Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:  "user does not exist",
			input: types.UpdatePasswordInput{ID: 3, Password: testPassword, ConfirmPassword: testPassword, NewPassword: "next-password"},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(3)).Return(&model.UserInfoModel{}, nil)
			},
			wantErr: ecode.ErrUserNotFound,
		},
		{
			name:  "current password is wrong",
			input: types.UpdatePasswordInput{ID: 3, Password: "wrong-password", ConfirmPassword: "wrong-password", NewPassword: "next-password"},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(3)).Return(stored, nil)
			},
			wantErr: ecode.ErrPasswordIncorrect,
		},
		{
			name:  "new password cannot be hashed",
			input: types.UpdatePasswordInput{ID: 3, Password: testPassword, ConfirmPassword: testPassword, NewPassword: longPassword},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(3)).Return(stored, nil)
			},
			wantErr: errcode.ErrEncrypt,
		},
		{
			name:  "update fails",
			input: types.UpdatePasswordInput{ID: 3, Password: testPassword, ConfirmPassword: testPassword, NewPassword: "next-password"},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(3)).Return(stored, nil)
				repo.EXPECT().UpdateUser(gomock.Any(), int64(3), gomock.Any()).Return(errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:  "success",
			input: types.UpdatePasswordInput{ID: 3, Password: testPassword, ConfirmPassword: testPassword, NewPassword: "next-password"},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(3)).Return(stored, nil)
				repo.EXPECT().UpdateUser(gomock.Any(), int64(3), gomock.Any()).DoAndReturn(
					func(_ context.Context, _ int64, data model.UserInfoModel) error {
						if !auth.ComparePasswords(data.Password, "next-password") {
							t.Error("new password was not hashed")
						}
						if data.UpdatedAt <= 0 {
							t.Error("updated_at was not set")
						}
						return nil
					},
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newMockedService(t)
			tt.setup(repo)

			got, err := svc.UpdatePassword(context.Background(), tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdatePassword() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdatePassword() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("UpdatePassword() = nil, want empty output")
			}
		})
	}
}

func TestUserServiceGetUser(t *testing.T) {
	errRepo := errors.New("repo failed")
	now := time.Now().Unix()
	stored := &model.UserInfoModel{
		ID:        4,
		Username:  testUsername,
		Email:     testEmail,
		Phone:     "10086",
		Nickname:  "ally",
		Status:    1,
		CreatedAt: now,
	}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockUserRepo)
		wantErr error
		wantNil bool
		wantID  int64
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(4)).Return(stored, nil)
			},
			wantID: stored.ID,
		},
		{
			name: "repository fails",
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(4)).Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name: "missing user converts to an empty result",
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUser(gomock.Any(), int64(4)).Return(nil, nil)
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newMockedService(t)
			tt.setup(repo)

			got, err := svc.GetUser(context.Background(), types.GetUserInput{ID: 4})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetUser() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetUser() unexpected error: %v", err)
			}
			if tt.wantNil {
				if got.User != nil {
					t.Fatalf("GetUser() user = %+v, want nil", got.User)
				}
				return
			}
			if got.User == nil || got.User.Id != tt.wantID || got.User.Username != testUsername || got.User.Email != testEmail {
				t.Fatalf("GetUser() = %+v", got.User)
			}
		})
	}
}

func TestUserServiceBatchGetUsers(t *testing.T) {
	errRepo := errors.New("repo failed")
	first := &model.UserInfoModel{ID: 1, Username: "one", Email: "one@example.com"}
	third := &model.UserInfoModel{ID: 3, Username: "three", Email: "three@example.com"}

	tests := []struct {
		name      string
		ids       []int64
		setup     func(repo *mocks.MockUserRepo)
		wantErr   error
		wantOrder []int64
	}{
		{
			name:    "empty ids",
			ids:     nil,
			setup:   func(repo *mocks.MockUserRepo) {},
			wantErr: ecode.ErrInvalidArgument,
		},
		{
			name: "repository fails",
			ids:  []int64{1},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().BatchGetUsers(gomock.Any(), []int64{1}).Return(nil, errRepo)
			},
			wantErr: ecode.ErrInternalError,
		},
		{
			name: "preserves request order and skips missing ids",
			ids:  []int64{1, 2, 3},
			setup: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().BatchGetUsers(gomock.Any(), []int64{1, 2, 3}).Return([]*model.UserInfoModel{third, first}, nil)
			},
			wantOrder: []int64{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newMockedService(t)
			tt.setup(repo)

			got, err := svc.BatchGetUsers(context.Background(), types.BatchGetUsersInput{IDs: tt.ids})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("BatchGetUsers() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("BatchGetUsers() unexpected error: %v", err)
			}
			if len(got.Users) != len(tt.wantOrder) {
				t.Fatalf("BatchGetUsers() len = %d, want %d", len(got.Users), len(tt.wantOrder))
			}
			for i, id := range tt.wantOrder {
				if got.Users[i].Id != id {
					t.Fatalf("user[%d].Id = %d, want %d", i, got.Users[i].Id, id)
				}
			}
		})
	}
}
