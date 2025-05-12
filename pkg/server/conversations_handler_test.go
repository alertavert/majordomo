// Author: M. Massenzio (marco@alertavert.com), 5/3/25

package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/conversations"
	"github.com/alertavert/gpt4-go/pkg/server"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/conversations endpoint", func() {
	var (
		router      *gin.Engine
		cfg         *config.Config
		assistant   *completions.Majordomo
		testThread  conversations.Thread
		projectName string
	)

	BeforeEach(func() {
		cfgLoc, err := MkTempConfigFile(TestConfigLocation)
		Expect(err).NotTo(HaveOccurred())

		cfg, err = config.LoadConfig(cfgLoc)
		Expect(err).NotTo(HaveOccurred())

		assistant, err = completions.NewMajordomo(cfg)
		Expect(err).NotTo(HaveOccurred())

		gin.SetMode(gin.TestMode)
		router = gin.New()
		server.SetupTestRoutes(router, assistant)

		// Use the first project from config for testing
		projectName = cfg.Projects[0].Name

		// Create a test thread
		testThread = conversations.Thread{
			ID:          "test-thread-id",
			Name:        "Test Thread",
			Assistant:   "test-assistant",
			Description: "Test thread description",
		}
		err = assistant.Threads.AddThread(projectName, testThread)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("GET /conversations/:thread_id", func() {
		Context("With valid parameters", func() {
			It("should return the specific thread", func() {
				req, _ := http.NewRequest("GET",
					fmt.Sprintf("/conversations/%s?project=%s", testThread.ID, projectName),
					nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))
				Expect(resp.Body.String()).To(ContainSubstring(testThread.ID))
				Expect(resp.Body.String()).To(ContainSubstring(testThread.Name))
				Expect(resp.Body.String()).To(ContainSubstring(testThread.Assistant))
			})
		})

		Context("With invalid parameters", func() {
			It("should return 400 when project parameter is missing", func() {
				req, _ := http.NewRequest("GET",
					fmt.Sprintf("/conversations/%s", testThread.ID),
					nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Body.String()).To(ContainSubstring("project query parameter is required"))
			})

			It("should return 404 when thread doesn't exist", func() {
				req, _ := http.NewRequest("GET",
					fmt.Sprintf("/conversations/nonexistent?project=%s", projectName),
					nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNotFound))
				Expect(resp.Body.String()).To(ContainSubstring("thread not found"))
			})
		})
	})
})
