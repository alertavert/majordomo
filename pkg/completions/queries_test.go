package completions_test

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	TestConfigLocation = "../../testdata/test_config.yaml"
)

var _ = Describe("Majordomo", func() {
	var (
		majordomo *completions.Majordomo
		cfg       *config.Config
		err       error
	)

	BeforeEach(func() {
		// Load configuration
		cfg, err = config.LoadConfig(TestConfigLocation)
		Expect(err).NotTo(HaveOccurred())

		// Create a new Majordomo instance
		majordomo, err = completions.NewMajordomo(cfg)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Session Handling", func() {
		It("can create a new session", func() {
			session := majordomo.NewSessionIfNotExists("new-session-id")
			Expect(session).NotTo(BeNil())
			Expect(session.SessionID).To(Equal("new-session-id"))
		})

		It("can retrieve an existing session", func() {
			majordomo.NewSessionIfNotExists("existing-session-id")
			session := majordomo.NewSessionIfNotExists("existing-session-id")
			Expect(session).NotTo(BeNil())
			Expect(session.SessionID).To(Equal("existing-session-id"))
		})

		It("can build messages for the session", func() {
			prompt := completions.PromptRequest{
				Prompt:   "What’s the weather like today?",
				Scenario: "some-scenario",
				Session:  "session-for-messages",
			}

			// We need to monkey-patch the scenario, since we don't have a real one
			completions.GetScenarios = func() *completions.Scenarios {
				return &completions.Scenarios{
					Common: "Here are some common instructions",
					Scenarios: map[string]string{
						"some-scenario": "Here are some instructions for the scenario",
					},
				}
			}

			// We need to initialize the session with the scenario first
			session := majordomo.NewSessionIfNotExists(prompt.Session)
			err := session.Init(prompt.Scenario)
			Expect(err).NotTo(HaveOccurred())

			messages, err := majordomo.BuildMessages(&prompt)
			Expect(err).NotTo(HaveOccurred())
			Expect(messages).NotTo(BeEmpty())
			Expect(len(messages)).To(Equal(3))
			Expect(messages[0].Content).To(Equal("Here are some common instructions"))
			Expect(messages[1].Content).To(Equal("Here are some instructions for the scenario"))
			Expect(messages[2].Content).To(Equal("What’s the weather like today?"))
		})
	})
})
