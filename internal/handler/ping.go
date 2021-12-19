package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/log"
)

// Ping ping
// @Summary ping
// @Description ping
// @Tags system
// @Accept  json
// @Produce  json
// @Router /ping [get]
func Ping(c *gin.Context) {
	log.Info("Get function called.")

	app.Success(c, gin.H{})
}
