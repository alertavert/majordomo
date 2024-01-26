package server_test

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/server"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)


var _ = Describe("/projects endpoint", func() {
	var (
		router    *gin.Engine
		cfg       *config.Config
		assistant *completions.Majordomo
	)

	BeforeEach(func() {
		cfgLoc, err := MkTempConfigFile(TestConfigLocation)
		Expect(err).NotTo(HaveOccurred())

		// Load configuration
		cfg, err = config.LoadConfig(cfgLoc)
		Expect(err).NotTo(HaveOccurred())

		// Create a new Majordomo instance
		assistant, err = completions.NewMajordomo(cfg)
		Expect(err).NotTo(HaveOccurred())

		// Set up the server
		gin.SetMode(gin.TestMode)
		router = gin.New()
		server.SetupTestRoutes(router, assistant)
	})

	Describe("GET /projects", func() {
		It("should return all the projects", func() {
			req, _ := http.NewRequest("GET", "/projects", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("should return the names of projects", func() {
			req, _ := http.NewRequest("GET", "/projects", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			body := resp.Body.String()
			for _, project := range cfg.Projects {
				Expect(body).To(ContainSubstring(project.Name))
			}
		})
	})

	Describe("GET /projects/:project_name", func() {
		Context("With a valid project name", func() {
			It("should return project details", func() {
				project := cfg.Projects[0]
				req, _ := http.NewRequest("GET", "/projects/"+project.Name, nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))
				Expect(resp.Body.String()).To(ContainSubstring(project.Name))
				Expect(resp.Body.String()).To(ContainSubstring(project.Description))
				Expect(resp.Body.String()).To(ContainSubstring(project.Location))
			})
		})

		Context("With an invalid project name", func() {
			It("should return a 404 error", func() {
				req, _ := http.NewRequest("GET", "/projects/nonexistent", nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("POST /projects", func() {
		Context("With valid project data", func() {
			It("should create a new project", func() {
				newProjectJson := `{"name":"new-project","description":"A new Project","location":"/some/path"}`
				req, _ := http.NewRequest("POST", "/projects", strings.NewReader(newProjectJson))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusCreated))

				// Check that new project was added to the configuration
				project := cfg.GetProject("new-project")
				Expect(project).NotTo(BeNil())
				Expect(project.Name).To(Equal("new-project"))
				Expect(project.Description).To(Equal("A new Project"))
				Expect(project.Location).To(Equal("/some/path"))
			})
		})

		Context("With invalid project data", func() {
			It("should return a 400 error", func() {
				newProjectJson := `{"Name": "", "Description": "A new Project", "Location":"/some/path" }`
				req, _ := http.NewRequest("POST", "/projects", strings.NewReader(newProjectJson))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("With a project name that already exists", func() {
			It("should return a 409 conflict error", func() {
				existingProjectName := cfg.Projects[0].Name
				newProjectJson := `{"name":"` + existingProjectName + `","description":"Duplicate Project","location":"/some/other/path"}`
				req, _ := http.NewRequest("POST", "/projects", strings.NewReader(newProjectJson))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusConflict))
			})
		})
	})

	Describe("PUT /projects/:project_name", func() {
		Context("With valid update data", func() {
			It("should update an existing project", func() {
				project := cfg.Projects[0]
				updateProjectJson := `{"name":"updated-name","description":"Updated Description","location":"/updated/path"}`
				req, _ := http.NewRequest("PUT", "/projects/"+project.Name, strings.NewReader(updateProjectJson))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))

				// Check that the project was updated
				project = *cfg.GetProject("updated-name")
				Expect(project).NotTo(BeNil())
				Expect(project.Description).To(Equal("Updated Description"))
				Expect(project.Location).To(Equal("/updated/path"))
			})
		})

		Context("With invalid update data", func() {
			It("should return a 400 error", func() {
				project := cfg.Projects[0]
				updateProjectJson := `{"name":""}`
				req, _ := http.NewRequest("PUT", "/projects/"+project.Name, strings.NewReader(updateProjectJson))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("With an invalid project name", func() {
			It("should return a 404 error", func() {
				updateProjectJson := `{"name":"updated-name","description":"Updated Description","location":"/updated/path"}`
				req, _ := http.NewRequest("PUT", "/projects/nonexistent", strings.NewReader(updateProjectJson))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("DELETE /projects/:project_name", func() {
		Context("With a valid project name", func() {
			It("should delete the project", func() {
				projectNameToDelete := cfg.Projects[0].Name
				req, _ := http.NewRequest("DELETE", "/projects/"+projectNameToDelete, nil)
				resp := httptest.NewRecorder()

				initialProjectCount := len(cfg.Projects)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))

				// Check that the project was deleted from the configuration
				project := cfg.GetProject(projectNameToDelete)
				Expect(project).To(BeNil())
				Expect(len(cfg.Projects)).To(Equal(initialProjectCount - 1))
			})
		})

		Context("With an invalid project name", func() {
			It("should return a 404 error", func() {
				req, _ := http.NewRequest("DELETE", "/projects/nonexistent", nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
