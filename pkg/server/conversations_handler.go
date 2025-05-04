// Author: M. Massenzio (marco@alertavert.com), 5/3/25

package server

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/conversations"
	"net/http"

	"github.com/gin-gonic/gin"
)

// threadsGetHandler handles GET requests to list all threads for a project
func threadsGetHandler(assistant *completions.Majordomo) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Query("project")
		if projectName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "project query parameter is required"})
			return
		}

		// Check if project exists
		project := assistant.Config.GetProject(projectName)
		if project == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("project '%s' not found", projectName)})
			return
		}
		threads := assistant.Threads.GetAllThreads(projectName)
		if threads == nil {
			threads = []conversations.Thread{}
		}

		c.JSON(http.StatusOK, gin.H{
			"project": projectName,
			"threads": threads,
		})
	}
}

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
