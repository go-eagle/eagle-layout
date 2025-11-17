package types

// User include user base info and user profile
type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginAt   int64  `json:"login_at"` // login time for last times
	Status    int32  `json:"status"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Gender    int32  `json:"gender"`
	Birthday  string `json:"birthday"`
	Bio       string `json:"bio"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// RegisterInput 注册输入
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// RegisterOutput 注册输出
type RegisterOutput struct {
	ID       int64
	Username string
}

// LoginInput 登录输入
type LoginInput struct {
	Email    string
	Username string
	Password string
}

// LoginOutput 登录输出
type LoginOutput struct {
	ID          int64
	AccessToken string
}

// CreateUserInput 创建用户输入
type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

// CreateUserOutput 创建用户输出
type CreateUserOutput struct {
	ID       int64
	Username string
}

// UpdateUserInput 更新用户输入
type UpdateUserInput struct {
	UserId   int64
	Nickname string
	Phone    string
	Email    string
	Avatar   string
	Gender   int32
	Birthday string
	Bio      string
	Status   int32
}

// UpdateUserOutput 更新用户输出
type UpdateUserOutput struct {
	UserId    int64
	Nickname  string
	Phone     string
	Email     string
	Avatar    string
	Gender    int32
	Birthday  string
	Bio       string
	Status    int32
	UpdatedAt int64
}

// UpdatePasswordInput 更新密码输入
type UpdatePasswordInput struct {
	ID              int64
	Password        string
	ConfirmPassword string
	NewPassword     string
}

// UpdatePasswordOutput 更新密码输出
type UpdatePasswordOutput struct{}

// GetUserInput 获取用户输入
type GetUserInput struct {
	ID int64
}

// GetUserOutput 获取用户输出
type GetUserOutput struct {
	User *User
}

// BatchGetUsersInput 批量获取用户输入
type BatchGetUsersInput struct {
	IDs []int64
}

// BatchGetUsersOutput 批量获取用户输出
type BatchGetUsersOutput struct {
	Users []*User
}

// LogoutInput 登出输入
type LogoutInput struct {
	ID          int64
	AccessToken string
}

// LogoutOutput 登出输出
type LogoutOutput struct{}
