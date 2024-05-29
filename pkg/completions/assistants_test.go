package completions_test

import (
	"github.com/alertavert/gpt4-go/pkg/completions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Assistants", func() {
	var (
		assistants *completions.Assistants
		err        error
	)

	Describe("ReadInstructions", func() {
		Context("when the YAML file is correctly formatted", func() {
			BeforeEach(func() {
				assistants, err = completions.ReadInstructions("../../testdata/test_assistants.yaml")
			})
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should correctly parse common instructions", func() {
				Expect(assistants.Common).To(Equal("common test scenario\n"))
			})
			It("should correctly parse individual instructions", func() {
				Expect(assistants.Instructions).To(HaveKeyWithValue("dev", ContainSubstring("You are an experienced Go developer;")))
				Expect(assistants.Instructions).To(HaveKeyWithValue("test", ContainSubstring("This is a test scenario")))
			})
		})
		Context("when the YAML file does not exist", func() {
			BeforeEach(func() {
				_, err = completions.ReadInstructions("../../testdata/invalid_path.yaml")
			})
			It("should return an error", func() {
   				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
		Context("when the YAML file is malformed", func() {
			BeforeEach(func() {
				assistants, err = completions.ReadInstructions("../../testdata/malformed_assistants.yaml")
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(assistants).To(BeNil())
			})
		})
	})

	Describe("GetInstructions", func() {
		BeforeEach(func() {
			assistants, err = completions.ReadInstructions("../../testdata/test_assistants.yaml")
			Expect(err).NotTo(HaveOccurred())
		})
		It("retrieves the correct instructions for a valid key", func() {
			Expect(assistants.GetInstructions("dev")).To(ContainSubstring("You are an experienced Go developer;"))
		})
		It("retrieves an empty string for a non-existent key", func() {
			Expect(assistants.GetInstructions("nonexistent")).To(Equal(""))
		})
	})

	Describe("Names", func() {
		BeforeEach(func() {
			assistants, err = completions.ReadInstructions("../../testdata/test_assistants.yaml")
			Expect(err).NotTo(HaveOccurred())
		})
		It("returns all the names of the configured assistants", func() {
			names := assistants.Names()
			Expect(names).To(ConsistOf([]string{"dev", "test"}))
		})
		It("returns an empty slice if no instructions are configured", func() {
			emptyAssistants := completions.Assistants{}
			Expect(emptyAssistants.Names()).To(BeEmpty())
		})
	})
})
