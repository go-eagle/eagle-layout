package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/go-eagle/eagle-layout/internal/dal/db/model"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	"github.com/go-eagle/eagle-layout/internal/types"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"

	pb "github.com/go-eagle/eagle-layout/api/user/v1"
)

// UserService 用户业务服务
type UserService struct {
	repo repository.UserRepo
}

// NewUserService 创建用户业务服务
func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, input types.RegisterInput) (*types.RegisterOutput, error) {
	// 检查邮箱是否存在
	userBase, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("[UserService] Register GetUserByEmail error: %w", err)
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, ecode.ErrUserIsExist
	}

	// 检查用户名是否存在
	userBase, err = s.repo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, ecode.ErrInternalError
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, ecode.ErrUserIsExist
	}

	// 生成密码哈希
	pwdHashed, err := auth.HashAndSalt(input.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// 创建用户
	user, err := s.newUser(input.Username, input.Email, pwdHashed)
	if err != nil {
		return nil, ecode.ErrInternalError
	}
	uid, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("[UserService] Register CreateUser error: %w", err)
	}

	// 发送欢迎邮件
	err = tasks.NewEmailWelcomeTask(tasks.EmailWelcomePayload{UserID: uid})
	if err != nil {
		return nil, fmt.Errorf("[UserService] Register NewEmailWelcomeTask error: %w", err)
	}

	return &types.RegisterOutput{
		ID:       uid,
		Username: input.Username,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, input types.LoginInput) (*types.LoginOutput, error) {
	var user *model.UserInfoModel
	var err error

	// 尝试用邮箱登录
	if input.Email != "" {
		user, err = s.repo.GetUserByEmail(ctx, input.Email)
		if err != nil {
			return nil, fmt.Errorf("[UserService] Login GetUserByEmail error: %w", err)
		}
	}

	// 如果没找到，尝试用户名
	if user == nil && input.Username != "" {
		user, err = s.repo.GetUserByUsername(ctx, input.Username)
		if err != nil {
			return nil, fmt.Errorf("[UserService] Login GetUserByUsername error: %w", err)
		}
	}

	// 用户不存在或密码错误
	if user == nil || user.ID == 0 || !auth.ComparePasswords(user.Password, input.Password) {
		return nil, ecode.ErrPasswordIncorrect
	}

	// 生成 token
	payload := map[string]interface{}{"user_id": user.ID, "username": user.Username}
	token, err := app.Sign(ctx, payload, app.Conf.JwtSecret, int64(cache.UserTokenExpireTime))
	if err != nil {
		return nil, ecode.ErrToken
	}

	// 记录 token 到 redis
	err = cache.NewUserTokenCache().SetUserTokenCache(ctx, user.ID, token, cache.UserTokenExpireTime)
	if err != nil {
		return nil, ecode.ErrToken
	}

	return &types.LoginOutput{
		ID:          user.ID,
		AccessToken: token,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, input types.LogoutInput) (*types.LogoutOutput, error) {
	c := cache.NewUserTokenCache()

	// 检查 token
	token, err := c.GetUserTokenCache(ctx, input.ID)
	if err != nil {
		return nil, ecode.ErrToken
	}
	if token != input.AccessToken {
		return nil, ecode.ErrAccessDenied
	}

	// 从缓存中删除 token
	err = c.DelUserTokenCache(ctx, input.ID)
	if err != nil {
		return nil, ecode.ErrInternalError
	}

	return &types.LogoutOutput{}, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, input types.CreateUserInput) (*types.CreateUserOutput, error) {
	// 生成密码哈希
	pwd, err := auth.HashAndSalt(input.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// 创建用户
	user, err := s.newUser(input.Username, input.Email, pwd)
	if err != nil {
		return nil, fmt.Errorf("[UserService] CreateUser newUser error: %w", err)
	}
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("[UserService] CreateUser CreateUser error: %w", err)
	}

	return &types.CreateUserOutput{
		ID:       id,
		Username: input.Username,
	}, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, input types.UpdateUserInput) (*types.UpdateUserOutput, error) {
	if input.UserId == 0 {
		return nil, ecode.ErrInvalidArgument
	}

	user := model.UserInfoModel{
		Nickname:  input.Nickname,
		Email:     input.Email,
		Avatar:    input.Avatar,
		Birthday:  input.Birthday,
		Bio:       input.Bio,
		Status:    input.Status,
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUser(ctx, input.UserId, user)
	if err != nil {
		return nil, fmt.Errorf("[UserService] UpdateUser UpdateUser error: %w", err)
	}

	return &types.UpdateUserOutput{
		UserId:    input.UserId,
		Nickname:  input.Nickname,
		Email:     input.Email,
		Avatar:    input.Avatar,
		Gender:    input.Gender,
		Birthday:  input.Birthday,
		Bio:       input.Bio,
		Status:    input.Status,
		UpdatedAt: time.Now().Unix(),
	}, nil
}

// UpdatePassword 更新密码
func (s *UserService) UpdatePassword(ctx context.Context, input types.UpdatePasswordInput) (*types.UpdatePasswordOutput, error) {
	if input.ID == 0 {
		return nil, ecode.ErrInvalidArgument
	}
	if input.Password == "" || input.ConfirmPassword == "" || input.NewPassword == "" {
		return nil, ecode.ErrInvalidArgument
	}
	if input.Password != input.ConfirmPassword {
		return nil, ecode.ErrTwicePasswordNotMatch
	}

	// 获取用户信息
	user, err := s.repo.GetUser(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("[UserService] UpdatePassword GetUser error: %w", err)
	}
	if user == nil || user.ID == 0 {
		return nil, ecode.ErrUserNotFound
	}

	// 验证旧密码
	if !auth.ComparePasswords(user.Password, input.Password) {
		return nil, ecode.ErrPasswordIncorrect
	}

	// 生成新密码哈希
	newPwd, err := auth.HashAndSalt(input.NewPassword)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// 更新密码
	data := model.UserInfoModel{
		Password:  newPwd,
		UpdatedAt: time.Now().Unix(),
	}
	err = s.repo.UpdateUser(ctx, user.ID, data)
	if err != nil {
		return nil, fmt.Errorf("[UserService] UpdatePassword UpdateUser error: %w", err)
	}

	return &types.UpdatePasswordOutput{}, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, input types.GetUserInput) (*types.GetUserOutput, error) {
	user, err := s.repo.GetUser(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("[UserService] GetUser GetUser error: %w", err)
	}

	u, err := s.convertUser(user)
	if err != nil {
		return nil, fmt.Errorf("[UserService] GetUser convertUser error: %w", err)
	}

	return &types.GetUserOutput{
		User: u,
	}, nil
}

// BatchGetUsers 批量获取用户
func (s *UserService) BatchGetUsers(ctx context.Context, input types.BatchGetUsersInput) (*types.BatchGetUsersOutput, error) {
	if len(input.IDs) == 0 {
		return nil, ecode.ErrInvalidArgument
	}

	// 获取用户列表
	userBases, err := s.repo.BatchGetUsers(ctx, input.IDs)
	if err != nil {
		return nil, ecode.ErrInternalError
	}

	userMap := make(map[int64]*model.UserInfoModel, 0)
	for _, val := range userBases {
		userMap[val.ID] = val
	}

	var users []*types.User
	for _, id := range input.IDs {
		user, ok := userMap[id]
		if !ok {
			continue
		}
		u, err := s.convertUser(user)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	return &types.BatchGetUsersOutput{
		Users: users,
	}, nil
}

// newUser 创建用户模型
func (s *UserService) newUser(username, email, password string) (model.UserInfoModel, error) {
	return model.UserInfoModel{
		Username:  username,
		Email:     email,
		Password:  password,
		Status:    int32(pb.StatusType_NORMAL),
		CreatedAt: time.Now().Unix(),
	}, nil
}

// convertUser 转换用户模型
func (s *UserService) convertUser(u *model.UserInfoModel) (*types.User, error) {
	if u == nil {
		return nil, nil
	}
	return &types.User{
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
	}, nil
}
