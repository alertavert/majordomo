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

	Describe("Majordomo", func() {
		It("can create a new Assistant", func() {
			Expect(majordomo).NotTo(BeNil())
			Expect(majordomo.Client).NotTo(BeNil())
			Expect(majordomo.Model).To(Equal(completions.DefaultModel))
		})
	})
})
