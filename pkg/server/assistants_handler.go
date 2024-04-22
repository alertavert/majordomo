package server

import (
	"context"
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// projectsGetHandler handles the GET request for the '/projects' endpoint.
func assistantsGetHandler(s *completions.Majordomo) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := s.Client.ListAssistants(context.Background(), nil, nil, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list.Assistants)
	}
}
