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

var _ = Describe("Parse Handler", func() {
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

	Describe("POST /parse", func() {
		Context("with valid request body", func() {
			It("should successfully parse a simple prompt", func() {
				promptReq := completions.PromptRequest{
					Prompt:    "Write a simple hello world program",
					Assistant: "default",
					ThreadName: "test-thread",
				}
				body, _ := json.Marshal(promptReq)
				req, _ := http.NewRequest("POST", "/parse", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusOK))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["response"]).To(Equal("success"))
				Expect(response["message"]).To(Equal(promptReq.Prompt))
			})

			It("should successfully parse a prompt with code references", func() {
				promptReq := completions.PromptRequest{
					Prompt:    "Update the code in `main.go`",
					Assistant: "default",
					ThreadName: "test-thread",
				}
				body, _ := json.Marshal(promptReq)
				req, _ := http.NewRequest("POST", "/parse", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusOK))
				var response map[string]interface{}
				Expect(json.Unmarshal(resp.Body.Bytes(), &response)).ShouldNot(HaveOccurred())
				Expect(response["response"]).To(Equal("success"))
			})
		})
		Context("with invalid request body", func() {
			It("should return 400 for missing prompt", func() {
				promptReq := map[string]string{
					"invalid_field": "some value",
				}
				body, _ := json.Marshal(promptReq)
				req, _ := http.NewRequest("POST", "/parse", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				Ω(json.Unmarshal(resp.Body.Bytes(), &response)).ToNot(HaveOccurred())
				Expect(response["response"]).To(Equal("error"))
			})

			It("should return 400 for malformed JSON", func() {
				malformedJSON := `{"prompt": "test", project": "test"}`
				req, _ := http.NewRequest("POST", "/parse", bytes.NewBuffer([]byte(malformedJSON)))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				Expect(response["response"]).To(Equal("error"))
			})

			It("should return 400 for empty request body", func() {
				req, _ := http.NewRequest("POST", "/parse", nil)
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				Expect(response["response"]).To(Equal("error"))
			})
		})
	})
})
