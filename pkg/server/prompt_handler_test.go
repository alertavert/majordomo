// Author: M. Massenzio (marco@alertavert.com), 5/3/25

package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/server"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prompt Handler", func() {
	var (
		router    *gin.Engine
		cfg       *config.Config
		assistant *completions.Majordomo
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
	})

	Describe("POST /prompt", func() {
		// TODO: we don't test "positive" cases here, as they are already covered in the completions
		// 		integration tests (completions/integration_test.go).
		// 		We should consider adding tests here using go-mocks.
		// TODO: Add tests to verify that thread_name is returned in the API response

		Context("with invalid request body", func() {
			It("should return 400 for missing prompt", func() {
				promptReq := map[string]string{
					"assistant": "default",
				}
				body, _ := json.Marshal(promptReq)
				req, _ := http.NewRequest("POST", "/prompt", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["status"]).To(Equal("error"))
			})

			It("should return 400 for missing assistant", func() {
				promptReq := map[string]string{
					"prompt": "Test prompt",
				}
				body, _ := json.Marshal(promptReq)
				req, _ := http.NewRequest("POST", "/prompt", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["status"]).To(Equal("error"))
			})

			It("should return 400 for malformed JSON", func() {
				malformedJSON := `{"prompt": "test", assistant": "default"}`
				req, _ := http.NewRequest("POST", "/prompt", bytes.NewBuffer([]byte(malformedJSON)))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["status"]).To(Equal("error"))
			})

			It("should return 400 for empty request body", func() {
				req, _ := http.NewRequest("POST", "/prompt", nil)
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["status"]).To(Equal("error"))
			})

			It("should return 400 when both thread ID and thread name are missing", func() {
				promptReq := map[string]string{
					"prompt":    "Test prompt",
					"assistant": "default",
					// Both thread_id and thread_name are intentionally omitted
				}
				body, _ := json.Marshal(promptReq)
				req, _ := http.NewRequest("POST", "/prompt", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["status"]).To(Equal("error"))
				// Verify the error message indicates the thread validation issue
				Expect(response["message"]).To(ContainSubstring("required"))
			})
		})
	})
})
