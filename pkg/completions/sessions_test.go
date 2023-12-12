package completions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sashabaranov/go-openai"

	"github.com/alertavert/gpt4-go/pkg/completions"
)

var _ = Describe("Session", func() {
	var (
		session *completions.Session
	)

	BeforeEach(func() {
		session = completions.NewSession("test-session-id")
	})

	Describe("Prompt and Response handling", func() {
		It("can add a prompt", func() {
			session.AddPrompt("Hello")
			Expect(len(session.GetUserPrompts())).To(Equal(1))
		})
		It("can add a response", func() {
			session.AddResponse("Hi there")
			Expect(len(session.GetBotResponses())).To(Equal(1))
		})
		It("knows if it's empty", func() {
			Expect(session.IsEmpty()).To(BeTrue())
			session.AddPrompt("Hi")
			Expect(session.IsEmpty()).To(BeFalse())
		})
		It("can get the conversation messages", func() {
			session.AddPrompt("Hi")
			session.AddResponse("Hello")
			conversation := session.GetConversation()
			Expect(conversation).To(HaveLen(2))

			// Ensure the order is correct: prompt then response
			Expect(conversation[0].Role).To(Equal(openai.ChatMessageRoleUser))
			Expect(conversation[1].Role).To(Equal(openai.ChatMessageRoleAssistant))
		})
	})

	Describe("Streamlining Conversations", func() {
		It("can clip the oldest messages", func() {
			session.AddPrompt("Hi")
			session.AddResponse("Hello")
			session.AddPrompt("How are you?")
			session.AddResponse("Good")

			session.Clip(1) // Removes the oldest prompt/response pair
			Expect(len(session.GetUserPrompts())).To(Equal(1))
			Expect(session.GetUserPrompts()[0].Content).To(Equal("How are you?"))
		})
		It("can clear the session", func() {
			session.AddPrompt("Hi")
			session.AddResponse("Hello")
			session.Clear()

			Expect(session.IsEmpty()).To(BeTrue())
		})
	})

	Describe("Session Initialization", func() {
		Context("with an initialized scenarios", func() {
			BeforeEach(func() {
				completions.GetScenarios = func() *completions.Scenarios {
					return &completions.Scenarios{
						Common: "Here are some common instructions",
						Scenarios: map[string]string{
							"mock-scenario-id": "Here are some instructions for the scenario",
							"another":          "Here are some other instructions for another scenario",
						},
					}
				}
			})
			It("should initialize a new session with scenario", func() {
				err := session.Init("mock-scenario-id")
				Ω(err).NotTo(HaveOccurred())
				Ω(session.ScenarioID).To(Equal("mock-scenario-id"))
				Ω(session.IsEmpty()).To(BeTrue())
				Ω(session.GetConversation()).Should(HaveLen(2))
			})
			It("should error when trying to re-initialize a session", func() {
				err := session.Init("mock-scenario-id")
				Expect(err).NotTo(HaveOccurred())
				err = session.Init("another")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("session test-session-id already has an ongoing conversation"))
			})
			It("should error when initializing a session with a non-existing scenario", func() {
				err := session.Init("non-existing-scenario")
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).To(ContainSubstring("no scenario found for non-existing-scenario"))
			})
		})
	})

	Describe("Error Conditions", func() {
		It("should not be possible to create a session with an empty ID", func() {
			newSession := completions.NewSession("")
			Expect(newSession).To(BeNil())
		})
		It("should handle creating a new session properly", func() {
			newSession := completions.NewSession("new-session-id")
			Expect(newSession).ToNot(BeNil())
			Expect(newSession.SessionID).To(Equal("new-session-id"))
		})
		It("should not create session for uninitialized scenarios", func() {
			err := session.Init("non-initialized-id")
			Expect(err).To(HaveOccurred())
			Ω(err.Error()).To(ContainSubstring("no scenarios found"))
		})
	})
})
