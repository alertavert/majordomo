package completions_test

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/preprocessors"
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
		Expect(majordomo).NotTo(BeNil())
	})

	Describe("Majordomo", func() {
		It("should return an error for an invalid API key", func() {
			id, err := majordomo.GetAssistantId("go_developer")
			Expect(err).To(HaveOccurred())
			Expect(id).To(BeEmpty())
		})
		It("should use the configured model", func() {
			cfg.Model = "gpt-4-turbo-preview"
			majordomo, err = completions.NewMajordomo(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(majordomo.Model).To(Equal("gpt-4-turbo-preview"))
		})
		It("will use the default model if not configured", func() {
			Expect(majordomo).NotTo(BeNil())
			Expect(majordomo.Client).NotTo(BeNil())
			Expect(majordomo.Model).To(Equal(completions.DefaultModel))
		})
		It("requires a valid active project to be configured", func() {
			cfg.ActiveProject = ""
			majordomo, err = completions.NewMajordomo(cfg)
			Expect(err).To(HaveOccurred())
			Expect(majordomo).To(BeNil())
		})
		It("will use the first project in the list as active, if none specified", func() {
			majordomo, err = completions.NewMajordomo(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(majordomo.Config.ActiveProject).To(Equal("test-project"))
		})
		It("can change the active project", func() {
			err := majordomo.SetActiveProject("test-project-2")
			Expect(err).NotTo(HaveOccurred())
			Expect(majordomo.Config.ActiveProject).To(Equal("test-project-2"))
		})
		It("will return an error if the project is not found", func() {
			err := majordomo.SetActiveProject("non-existent-project")
			Expect(err).To(HaveOccurred())
		})
		It("Will use the correct location for the code snippets", func() {
			// First check that the Config has the correct values
			cfg := majordomo.Config
			prj := cfg.GetActiveProject()
			Expect(prj).NotTo(BeNil())
			Expect(prj.ResolvedCodeSnippetsDir).To(HaveSuffix("test/location/.majordomo"))
			// This should also be what the CodeStore uses
			Expect(majordomo.CodeStore).NotTo(BeNil())
			// We cast the CodeStore to be a FilesystemStore to access the Location field.
			fsStore := majordomo.CodeStore.(*preprocessors.FilesystemStore)
			Expect(fsStore.DestCodeDir).To(HaveSuffix("test/location/.majordomo"))
		})
		It("will set the CodeStore to the new project's location", func() {
			err := majordomo.SetActiveProject("test-project-2")
			Expect(err).NotTo(HaveOccurred())
			Expect(majordomo.CodeStore).NotTo(BeNil())
			// We cast the CodeStore to be a FilesystemStore to access the Location field.
			fsStore := majordomo.CodeStore.(*preprocessors.FilesystemStore)
			Expect(fsStore.SourceCodeDir).To(Equal("test/location-2"))
			// The destination for the code returned by the bot should be as
			// configured in the test_config.yaml file, ending with the project name.
			Expect(fsStore.DestCodeDir).To(HaveSuffix("test/location-2/.majordomo"))
		})
	})
	Describe("When parsing a user prompt", func() {
		It("should successfully fill in the correct content from the source code map", func() {
			Expect(majordomo.SetActiveProject("actual")).NotTo(HaveOccurred())
			prompt := "Please update this code:\n'''sample/main.go\n" +
				"'''to also print the current date."
			request := completions.PromptRequest{
				Assistant: "go_developer",
				ThreadId:  "",
				Prompt:    prompt,
			}
			Expect(majordomo.PreparePrompt(&request)).ShouldNot(HaveOccurred())
			// Read the contents of the file from the filesystem
			// Remember that the code snippets are stored in the SourceCodeDir relative
			// to the project's location.
			code := &preprocessors.SourceCodeMap{
				"sample/main.go": "",
			}
			Expect(majordomo.CodeStore.GetSourceCode(code)).NotTo(HaveOccurred())
			contents, found := (*code)["sample/main.go"]
			Expect(found).To(BeTrue())
			Expect(request.Prompt).To(ContainSubstring(contents))
		})
		It("should fail for an invalid file path", func() {
			err := majordomo.SetActiveProject("actual")
			Expect(err).NotTo(HaveOccurred())

			prompt := "Please update this code:\n'''invalid/file/path\n'''"
			request := completions.PromptRequest{
				Assistant: "go_developer",
				ThreadId:  "",
				Prompt:    prompt,
			}
			Expect(majordomo.PreparePrompt(&request)).To(HaveOccurred())
		})
		It("can import multiple files", func() {
			err := majordomo.SetActiveProject("actual")
			Expect(err).NotTo(HaveOccurred())

			prompt := "Please update this code:\n'''sample/main.go\n" +
				"'''to also print the current date.\n" +
				"'''pkg/simple.go\n'''"
			request := completions.PromptRequest{
				Assistant: "go_developer",
				ThreadId:  "",
				Prompt:    prompt,
			}
			Expect(majordomo.PreparePrompt(&request)).ShouldNot(HaveOccurred())
			// Read the contents of the files from the filesystem
			code := &preprocessors.SourceCodeMap{
				"sample/main.go": "",
				"pkg/simple.go":  "",
			}
			contents, found := (*code)["sample/main.go"]
			Expect(found).To(BeTrue())
			Expect(request.Prompt).To(ContainSubstring(contents))
			contents, found = (*code)["pkg/simple.go"]
			Expect(found).To(BeTrue())
			Expect(request.Prompt).To(ContainSubstring(contents))
		})
	})
	Describe("When processing a prompt", func() {
		It("should fail to create a new thread if the API key is invalid", func() {
			cfg.OpenAIApiKey = "invalid"
			m, err := completions.NewMajordomo(cfg)
			Expect(err).NotTo(HaveOccurred())
			tid := m.CreateNewThread("My Project", "go_developer")
			Expect(tid).To(BeEmpty())
		})
		It("should return an error if the project is not found", func() {
			tid := majordomo.CreateNewThread("non-existent-project", "go_developer")
			Expect(tid).To(BeEmpty())
		})
	})
})
