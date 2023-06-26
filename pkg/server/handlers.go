package server

import (
	"github.com/gin-gonic/gin"
)

func setupHandlers(r *gin.Engine) {
	r.GET("/", indexHandler)
	r.POST("/prompt", promptHandler)
}
