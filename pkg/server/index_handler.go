package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Please use the /prompt?scenario=<scenario> endpoint to send a prompt to the server",
	})
}
