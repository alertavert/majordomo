// Author: M. Massenzio (marco@alertavert.com), 5/3/25

package server

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"net/http"

	"github.com/gin-gonic/gin"
)

// threadGetByIdHandler handles GET requests for a specific thread
func threadGetByIdHandler(assistant *completions.Majordomo) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Query("project")
		threadId := c.Param("thread_id")

		if projectName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "project query parameter is required"})
			return
		}

		thread, found := assistant.Threads.GetThread(projectName, threadId)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
			return
		}

		c.JSON(http.StatusOK, thread)
	}
}
