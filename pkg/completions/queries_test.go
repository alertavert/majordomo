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
		It("can create a new Assistant with the default model", func() {
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
	It("will set the CodeStore to the new project's location", func() {
		err := majordomo.SetActiveProject("test-project-2")
		Expect(err).NotTo(HaveOccurred())
		Expect(majordomo.CodeStore).NotTo(BeNil())
		// We cast the CodeStore to be a FilesystemStore to access the Location field.
		fsStore := majordomo.CodeStore.(*preprocessors.FilesystemStore)
		Expect(fsStore.SourceCodeDir).To(Equal("test/location-2"))
		// The destination for the code returned by the bot should be as
		// configured in the test_config.yaml file, ending with the project name.
		Expect(fsStore.DestCodeDir).To(HaveSuffix("code/snippets/test-project-2"))
	})
})
