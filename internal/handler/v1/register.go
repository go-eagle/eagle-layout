package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/errcode"
)

// RegisterHandler 包含 UserService
type RegisterHandler struct {
	UserService service.UserService
}

// NewRegisterHandler 创建一个新的 RegisterHandler
func NewRegisterHandler(userService service.UserService) *RegisterHandler {
	return &RegisterHandler{UserService: userService}
}

// RegisterHandler 登录
// @Summary 用户名和密码登录
// @Description demo
// @Tags user
// @Accept  json
// @Produce  json
// @Router /v1/auth/Register [post]
func (h *RegisterHandler) RegisterHandler(c *gin.Context) {
	// 从请求中提取参数（例如 JSON 或表单参数）
	var RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&RegisterRequest); err != nil {
		app.Error(c, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// 调用 UserService 的逻辑
	user, err := h.UserService.Register(c.Request.Context(), RegisterRequest.Username, RegisterRequest.Password)
	if err != nil {
		app.Error(c, errcode.ErrUnauthorized.WithDetails(err.Error()))
		return
	}

	// 返回成功响应
	app.Success(c, gin.H{"user": user})
}
