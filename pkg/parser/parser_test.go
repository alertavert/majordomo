/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package parser_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alertavert/gpt4-go/pkg/parser"
)

const (
	f1 = `package testdata

func test(s string) error {
	return nil
}
`
	f2 = `package testdata

import (
	"io"
	"os"
)

func anoter(f *os.File) ([]byte, error) {
	return io.ReadAll(f)
}
`
)

var _ = Describe("ParseBotResponse", func() {
	It("should return an error when no valid code snippets are found", func() {
		sourceCode, err := parser.ParseBotResponse("This is not a code snippet")
		Expect(err).To(HaveOccurred())
		Expect(sourceCode).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
	})

	Context("with a well-formed code snippet", func() {

		var source parser.SourceCode
		var err error

		BeforeEach(func() {
			source, err = parser.ParseBotResponse("'''test.txt\nThis is a test content\n'''")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully extract the correct content to the source code map", func() {
			expectedContent := "This is a test content\n"
			Expect(len(source)).To(Equal(1))
			Expect(source["test.txt"]).To(Equal(expectedContent))
		})
	})

	Context("with malformed code snippets", func() {
		It("should return an error when missing closing triple-quotes", func() {
			source, err := parser.ParseBotResponse("'''test.txt\n missing closing triple-quotes")
			Expect(err).To(HaveOccurred())
			Expect(source).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})

		It("should return an error when missing opening triple-quotes", func() {
			source, err := parser.ParseBotResponse("missing opening triple-quotes\n'''")
			Expect(err).To(HaveOccurred())
			Expect(source).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})

		It("should return an error when the file path is malformed", func() {
			source, err := parser.ParseBotResponse("'''\nthis/is not: a /valid]path\nThis is content\n'''")
			Expect(err).To(HaveOccurred())
			Expect(source).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no valid code snippets found"))
		})
	})
})
var _ = Describe("Parser", func() {
	Context("InsertSourceCode", func() {
		It("Should return the right file content", func() {
			text := "'''../../testdata/test1.go\n'''"
			expectedResult := fmt.Sprintf("'''../../testdata/test1.go\n%s'''", f1)
			result, err := parser.InsertSourceCode(text)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedResult))
		})

		It("Should return error if file not found", func() {
			text := "'''path/to/file.go\n'''"
			_, err := parser.InsertSourceCode(text)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(fmt.Sprintf(parser.ErrorNoCodeSnippetsFound,
				"path/to/file.go", "")))
		})
		It("Should insert multiple snippets", func() {
			text := `Some intro text
'''../../testdata/test1.go
%s'''
Some random text:
'''../../testdata/test2.go
%s'''
and some more text.`
			result, err := parser.InsertSourceCode(fmt.Sprintf(text, "", ""))
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(fmt.Sprintf(text, f1, f2)))
		})
	})
})
