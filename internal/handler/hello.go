package handler

import (
	"net/http"

	"github.com/go-eagle/eagle-layout/internal/service"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/app"
)

// Param 请求参数
type Param struct {
	Name string `form:"name"`
}

// Hello a demo handler
func Hello(c *gin.Context) {
	var p Param
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "invalid param",
		})
		return
	}
	res, err := service.UserSvc.Hello(c, p.Name)
	if err != nil {
		app.Error(c, err)
		return
	}
	app.Success(c, gin.H{
		"result": res,
	})
}
