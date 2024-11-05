package service

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	pb "github.com/go-eagle/eagle-layout/api/user/v1"
	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	"github.com/go-eagle/eagle-layout/internal/types"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"
)

var (
	_ pb.UserServiceServer = (*UserServiceServer)(nil)
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	repo repository.UserRepo
}

func NewUserServiceServer(repo repository.UserRepo) *UserServiceServer {
	return &UserServiceServer{
		repo: repo,
	}
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	err := req.Validate()
	if err != nil {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	var userBase *model.UserInfoModel
	// check user is existed
	userBase, err = s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, ecode.ErrUserIsExist.Status(req).Err()
	}
	userBase, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, ecode.ErrUserIsExist.Status(req).Err()
	}

	// gen a hash password
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// create a new user
	user, err := newUser(req.Username, req.Email, pwd)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	uid, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	// send welcome email
	err = tasks.NewEmailWelcomeTask(tasks.EmailWelcomePayload{UserID: uid})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.RegisterReply{
		Id:       uid,
		Username: req.Username,
	}, nil
}

func newUser(username, email, password string) (model.UserInfoModel, error) {
	return model.UserInfoModel{
		Username:  username,
		Email:     email,
		Password:  password,
		Status:    int32(pb.StatusType_NORMAL),
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	if len(req.Email) == 0 && len(req.Username) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}

	// get user base info
	var (
		user *model.UserInfoModel
		err  error
	)
	if req.Email != "" {
		user, err = s.repo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if user == nil && len(req.Username) > 0 {
		user, err = s.repo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if user == nil || user.ID == 0 {
		return nil, ecode.ErrPasswordIncorrect.Status(req).Err()
	}

	if !auth.ComparePasswords(user.Password, req.Password) {
		return nil, ecode.ErrPasswordIncorrect.Status(req).Err()
	}

	// Sign the json web token.
	payload := map[string]interface{}{"user_id": user.ID, "username": user.Username}
	token, err := app.Sign(ctx, payload, app.Conf.JwtSecret, int64(cache.UserTokenExpireTime))
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}

	// record token to redis
	err = cache.NewUserTokenCache().SetUserTokenCache(ctx, user.ID, token, cache.UserTokenExpireTime)
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}

	return &pb.LoginReply{
		Id:          user.ID,
		AccessToken: token,
	}, nil
}

func (s *UserServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	c := cache.NewUserTokenCache()
	// check token
	token, err := c.GetUserTokenCache(ctx, req.Id)
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}
	if token != req.AccessToken {
		return nil, ecode.ErrAccessDenied.Status(req).Err()
	}

	// delete token from cache
	err = c.DelUserTokenCache(ctx, req.GetId())
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.LogoutReply{}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	// gen a hash password
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// create a new user
	user, err := newUser(req.Username, req.Email, pwd)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.CreateUserReply{
		Id:       id,
		Username: req.Username,
	}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	if req.UserId == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}

	user := model.UserInfoModel{
		Nickname: req.Nickname,
		//Phone:     req.Phone,
		Email:  req.Email,
		Avatar: req.Avatar,
		//Gender:    cast.ToString(req.Gender),
		Birthday:  req.Birthday,
		Bio:       req.Bio,
		Status:    cast.ToInt32(req.Status),
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUser(ctx, req.UserId, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.UpdateUserReply{
		UserId:    req.UserId,
		Nickname:  req.Nickname,
		Phone:     req.Phone,
		Email:     req.Email,
		Avatar:    req.Avatar,
		Gender:    req.Gender,
		Birthday:  req.Birthday,
		Bio:       req.Bio,
		Status:    req.Status,
		UpdatedAt: time.Now().Unix(),
	}, nil
}

func (s *UserServiceServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	if len(req.Id) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}
	if len(req.Password) == 0 || len(req.NewPassword) == 0 || len(req.ConfirmPassword) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}
	if req.NewPassword != req.ConfirmPassword {
		return nil, ecode.ErrTwicePasswordNotMatch.Status(req).Err()
	}

	// get user base info
	var (
		user *model.UserInfoModel
		err  error
	)
	user, err = s.repo.GetUser(ctx, cast.ToInt64(req.Id))
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if user == nil || user.ID == 0 {
		return nil, ecode.ErrUserNotFound.Status(req).Err()
	}

	if !auth.ComparePasswords(user.Password, req.Password) {
		return nil, ecode.ErrPasswordIncorrect.Status(req).Err()
	}

	newPwd, err := auth.HashAndSalt(req.NewPassword)
	if err != nil {
		return nil, ecode.ErrEncrypt.Status(req).Err()
	}

	data := model.UserInfoModel{
		Password:  newPwd,
		UpdatedAt: time.Now().Unix(),
	}
	err = s.repo.UpdateUser(ctx, user.ID, data)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.UpdatePasswordReply{}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.repo.GetUser(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrUserNotFound.Status(req).Err()
		}
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	u, err := convertUser(user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.GetUserReply{
		User: u,
	}, nil
}

func (s *UserServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersReply, error) {
	// check rpc request if canceled
	if ctx.Err() == context.Canceled {
		return nil, ecode.ErrCanceled.Status(req).Err()
	}

	if len(req.GetIds()) == 0 {
		return nil, errors.New("ids is empty")
	}
	var (
		ids   []int64
		users []*pb.User
	)
	ids = req.GetIds()

	// user base
	userBases, err := s.repo.BatchGetUsers(ctx, ids)
	if err != nil {
		return nil, ecode.ErrInternalError.Status(req).Err()
	}
	userMap := make(map[int64]*model.UserInfoModel, 0)
	for _, val := range userBases {
		userMap[val.ID] = val
	}

	// compose data
	for _, id := range ids {
		user, ok := userMap[id]
		if !ok {
			continue
		}
		u, err := convertUser(user)
		if err != nil {
			// record log
			continue
		}
		users = append(users, u)
	}

	return &pb.BatchGetUsersReply{
		Users: users,
	}, nil
}

func convertUser(u *model.UserInfoModel) (*pb.User, error) {
	if u == nil {
		return nil, nil
	}
	user := &types.User{
		Id:        u.ID,
		Username:  u.Username,
		Phone:     u.Phone,
		Email:     u.Email,
		LoginAt:   u.LoginAt,
		Status:    u.Status,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	// copy to pb.user
	pbUser := &pb.User{}
	err := copier.Copy(pbUser, &user)
	if err != nil {
		return nil, err
	}
	return pbUser, nil
}
