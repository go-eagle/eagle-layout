package service

import (
	"context"
	"errors"
	"testing"

	"github.com/go-eagle/eagle/pkg/errcode"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	pb "github.com/go-eagle/eagle-layout/api/user/v1"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/mocks"
	"github.com/go-eagle/eagle-layout/internal/types"
)

func TestConvertToGrpcError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		code    codes.Code
		message string
	}{
		{
			name:    "record not found",
			err:     gorm.ErrRecordNotFound,
			code:    codes.NotFound,
			message: "user not found",
		},
		{
			name:    "invalid argument",
			err:     ecode.ErrInvalidArgument,
			code:    codes.InvalidArgument,
			message: errcode.ErrInvalidParam.Msg(),
		},
		{
			name:    "unknown business code",
			err:     ecode.ErrPasswordIncorrect,
			code:    codes.Unknown,
			message: ecode.ErrPasswordIncorrect.Msg(),
		},
		{
			name:    "plain error",
			err:     errors.New("boom"),
			code:    codes.Internal,
			message: "internal error: boom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertToGrpcError(tt.err)
			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("convertToGrpcError() = %v, want a status", err)
			}
			if st.Code() != tt.code || st.Message() != tt.message {
				t.Fatalf("status = %s %q, want %s %q", st.Code(), st.Message(), tt.code, tt.message)
			}
		})
	}
}

func TestUserServiceServerRegisterAndLogin(t *testing.T) {
	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		call    func(srv *UserServiceServer) error
		setup   func(svc *mocks.MockUserService)
		wantErr codes.Code
	}{
		{
			name: "register success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Register(gomock.Any(), types.RegisterInput{
					Username: testUsername, Email: testEmail, Password: testPassword,
				}).Return(&types.RegisterOutput{ID: 8, Username: testUsername}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.Register(context.Background(), &pb.RegisterRequest{
					Username: testUsername, Email: testEmail, Password: testPassword,
				})
				if err != nil {
					return err
				}
				if reply.Id != 8 || reply.Username != testUsername {
					return errors.New("unexpected register reply")
				}
				return nil
			},
		},
		{
			name: "register failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, ecode.ErrUserIsExist)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.Register(context.Background(), &pb.RegisterRequest{Username: testUsername})
				return err
			},
			wantErr: codes.Unknown,
		},
		{
			name: "login success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Login(gomock.Any(), types.LoginInput{
					Email: testEmail, Username: testUsername, Password: testPassword,
				}).Return(&types.LoginOutput{ID: 8, AccessToken: "tok"}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.Login(context.Background(), &pb.LoginRequest{
					Email: testEmail, Username: testUsername, Password: testPassword,
				})
				if err != nil {
					return err
				}
				if reply.Id != 8 || reply.AccessToken != "tok" {
					return errors.New("unexpected login reply")
				}
				return nil
			},
		},
		{
			name: "login failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errBoom)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.Login(context.Background(), &pb.LoginRequest{})
				return err
			},
			wantErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockUserService(ctrl)
			tt.setup(svc)
			err := tt.call(NewUserServiceServer(svc))
			assertGRPCCode(t, err, tt.wantErr)
		})
	}
}

func TestUserServiceServerProfile(t *testing.T) {
	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		setup   func(svc *mocks.MockUserService)
		call    func(srv *UserServiceServer) error
		wantErr codes.Code
	}{
		{
			name: "logout success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Logout(gomock.Any(), types.LogoutInput{ID: 3, AccessToken: "tok"}).
					Return(&types.LogoutOutput{}, nil)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.Logout(context.Background(), &pb.LogoutRequest{Id: 3, AccessToken: "tok"})
				return err
			},
		},
		{
			name: "logout failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(nil, ecode.ErrAccessDenied)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.Logout(context.Background(), &pb.LogoutRequest{})
				return err
			},
			wantErr: codes.PermissionDenied,
		},
		{
			name: "create user success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().CreateUser(gomock.Any(), types.CreateUserInput{
					Username: testUsername, Email: testEmail, Password: testPassword,
				}).Return(&types.CreateUserOutput{ID: 4, Username: testUsername}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
					Username: testUsername, Email: testEmail, Password: testPassword,
				})
				if err != nil {
					return err
				}
				if reply.Id != 4 || reply.Username != testUsername {
					return errors.New("unexpected create reply")
				}
				return nil
			},
		},
		{
			name: "create user failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errBoom)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{})
				return err
			},
			wantErr: codes.Internal,
		},
		{
			name: "update user maps enums",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().UpdateUser(gomock.Any(), types.UpdateUserInput{
					UserId: 5, Nickname: "ally", Phone: "10086", Email: testEmail,
					Avatar: "a.png", Gender: 1, Birthday: "2000-01-01", Bio: "hi", Status: 2,
				}).Return(&types.UpdateUserOutput{
					UserId: 5, Nickname: "ally", Phone: "10086", Email: testEmail,
					Avatar: "a.png", Gender: 1, Birthday: "2000-01-01", Bio: "hi", Status: 2, UpdatedAt: 99,
				}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.UpdateUser(context.Background(), &pb.UpdateUserRequest{
					UserId: 5, Nickname: "ally", Phone: "10086", Email: testEmail,
					Avatar: "a.png", Gender: pb.GenderType(1), Birthday: "2000-01-01", Bio: "hi", Status: pb.StatusType(2),
				})
				if err != nil {
					return err
				}
				if reply.UserId != 5 || reply.Gender != pb.GenderType(1) || reply.Status != pb.StatusType(2) || reply.UpdatedAt != 99 {
					return errors.New("unexpected update reply")
				}
				return nil
			},
		},
		{
			name: "update user failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil, errBoom)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.UpdateUser(context.Background(), &pb.UpdateUserRequest{})
				return err
			},
			wantErr: codes.Internal,
		},
		{
			name:  "update password rejects a bad id",
			setup: func(svc *mocks.MockUserService) {},
			call: func(srv *UserServiceServer) error {
				_, err := srv.UpdatePassword(context.Background(), &pb.UpdatePasswordRequest{Id: "abc"})
				return err
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "update password success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().UpdatePassword(gomock.Any(), types.UpdatePasswordInput{
					ID: 6, Password: "old", NewPassword: "new", ConfirmPassword: "old",
				}).Return(&types.UpdatePasswordOutput{}, nil)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.UpdatePassword(context.Background(), &pb.UpdatePasswordRequest{
					Id: "6", Password: "old", NewPassword: "new", ConfirmPassword: "old",
				})
				return err
			},
		},
		{
			name: "update password failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(nil, ecode.ErrUserNotFound)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.UpdatePassword(context.Background(), &pb.UpdatePasswordRequest{Id: "6"})
				return err
			},
			wantErr: codes.Unknown,
		},
		{
			name: "get user success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().GetUser(gomock.Any(), types.GetUserInput{ID: 4}).
					Return(&types.GetUserOutput{User: &types.User{Id: 4, Username: testUsername, Email: testEmail}}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.GetUser(context.Background(), &pb.GetUserRequest{Id: 4})
				if err != nil {
					return err
				}
				if reply.User == nil || reply.User.Id != 4 || reply.User.Username != testUsername {
					return errors.New("unexpected get user reply")
				}
				return nil
			},
		},
		{
			name: "get user with empty profile",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&types.GetUserOutput{}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.GetUser(context.Background(), &pb.GetUserRequest{Id: 4})
				if err != nil {
					return err
				}
				if reply.User != nil {
					return errors.New("expected nil user")
				}
				return nil
			},
		},
		{
			name: "get user failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errBoom)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.GetUser(context.Background(), &pb.GetUserRequest{})
				return err
			},
			wantErr: codes.Internal,
		},
		{
			name:  "batch get canceled",
			setup: func(svc *mocks.MockUserService) {},
			call: func(srv *UserServiceServer) error {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := srv.BatchGetUsers(ctx, &pb.BatchGetUsersRequest{Ids: []int64{1}})
				return err
			},
			wantErr: codes.Canceled,
		},
		{
			name: "batch get success",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().BatchGetUsers(gomock.Any(), types.BatchGetUsersInput{IDs: []int64{1, 2}}).
					Return(&types.BatchGetUsersOutput{Users: []*types.User{{Id: 1, Username: "one"}, {Id: 2, Username: "two"}}}, nil)
			},
			call: func(srv *UserServiceServer) error {
				reply, err := srv.BatchGetUsers(context.Background(), &pb.BatchGetUsersRequest{Ids: []int64{1, 2}})
				if err != nil {
					return err
				}
				if len(reply.Users) != 2 || reply.Users[0].Username != "one" || reply.Users[1].Id != 2 {
					return errors.New("unexpected batch reply")
				}
				return nil
			},
		},
		{
			name: "batch get failure",
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().BatchGetUsers(gomock.Any(), gomock.Any()).Return(nil, ecode.ErrInternalError)
			},
			call: func(srv *UserServiceServer) error {
				_, err := srv.BatchGetUsers(context.Background(), &pb.BatchGetUsersRequest{Ids: []int64{1}})
				return err
			},
			wantErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockUserService(ctrl)
			tt.setup(svc)
			err := tt.call(NewUserServiceServer(svc))
			assertGRPCCode(t, err, tt.wantErr)
		})
	}
}

func assertGRPCCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if want == codes.OK {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != want {
		t.Fatalf("error = %v, want code %s", err, want)
	}
}
