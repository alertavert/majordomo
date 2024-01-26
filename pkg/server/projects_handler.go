package server

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// projectsGetHandler handles the GET request for the '/projects' endpoint.
func projectsGetHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectNames := make([]string, len(cfg.Projects))
		for i, project := range cfg.Projects {
			projectNames[i] = project.Name
		}
		c.JSON(http.StatusOK, projectNames)
	}
}

// projectDetailsGetHandler handles the GET request for the '/projects/:project_name' endpoint.
func projectDetailsGetHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Param("project_name")
		project := cfg.GetProject(projectName)
		if project == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusOK, project)
	}
}

// Helper function to update project fields if the provided data is non-empty.
func updateProjectIfNotEmpty(original *config.Project, updates config.Project) {
	if updates.Name != "" && isProjectNameValid(updates.Name) {
		original.Name = updates.Name
	}
	if updates.Location != "" {
		original.Location = updates.Location
	}
	if updates.Description != "" {
		original.Description = updates.Description
	}
}

// Helper function to check if project name is valid
func isProjectNameValid(name string) bool {
	return len(name) > 0 && !strings.ContainsAny(name, " /?%#")
}

// Helper function to check if project name already exists
func isProjectNameExists(name string, projects []config.Project) bool {
	for _, project := range projects {
		if project.Name == name {
			return true
		}
	}
	return false
}

// projectPutHandler handles the PUT request for the '/projects/:project_name' endpoint.
func projectPutHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Param("project_name")

		var updatedProject config.Project
		if err := c.BindJSON(&updatedProject); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		// Check if the projectName is valid.
		if !isProjectNameValid(updatedProject.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project name"})
			return
		}

		projectIndex := -1
		for i, project := range cfg.Projects {
			if project.Name == projectName {
				projectIndex = i
				break
			}
		}

		if projectIndex == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		// Update only the fields that have been provided in the request body.
		updateProjectIfNotEmpty(&cfg.Projects[projectIndex], updatedProject)

		if err := cfg.Save(""); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
			return
		}
		c.JSON(http.StatusOK, cfg.Projects[projectIndex])
	}
}

// projectPostHandler handles the POST request for the '/projects' endpoint.
func projectPostHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProject config.Project
		if err := c.BindJSON(&newProject); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		fmt.Println(newProject)

		// Check for the validity of the project name and uniqueness.
		if !isProjectNameValid(newProject.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project name contains invalid characters"})
			return
		}
		if isProjectNameExists(newProject.Name, cfg.Projects) {
			c.JSON(http.StatusConflict, gin.H{"error": "Project already exists"})
			return
		}

		cfg.Projects = append(cfg.Projects, newProject)
		if err := cfg.Save(""); err != nil {
			errMsg := fmt.Sprintf("Failed to save new project: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusCreated, newProject)
	}
}

// projectDeleteHandler handles the DELETE request for the '/projects/:project_name' endpoint.
func projectDeleteHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Param("project_name")
		projectIndex := -1
		for i, project := range cfg.Projects {
			if project.Name == projectName {
				projectIndex = i
				break
			}
		}
		if projectIndex == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		cfg.Projects = append(cfg.Projects[:projectIndex], cfg.Projects[projectIndex+1:]...)
		if err := cfg.Save(""); err != nil {
			errMsg := fmt.Sprintf("Failed to delete project: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
	}
}
