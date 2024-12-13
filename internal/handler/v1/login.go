package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"
)

// LoginHandler 包含 UserService
type LoginHandler struct {
	UserService service.UserService
}

// NewLoginHandler 创建一个新的 LoginHandler
func NewLoginHandler(userService service.UserService) *LoginHandler {
	return &LoginHandler{UserService: userService}
}

// LoginHandler 登录
// @Summary 用户名和密码登录
// @Description demo
// @Tags user
// @Accept  json
// @Produce  json
// @Router /v1/auth/login [post]
func (h *LoginHandler) LoginHandler(c *gin.Context) {
	// 从请求中提取参数（例如 JSON 或表单参数）
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// 调用 UserService 的逻辑
	user, err := h.UserService.Login(c.Request.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		app.Error(c, errcode.ErrUnauthorized.WithDetails(err.Error()))
		return
	}

	// 返回成功响应
	app.Success(c, gin.H{"user": user})
}
