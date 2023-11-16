/*
 * Copyright (c) 2023 AlertAvert.com. All rights reserved.
 */

package preprocessors_test

import (
	"fmt"
	"github.com/alertavert/gpt4-go/pkg/preprocessors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var f1 = `func promptHandler(c *gin.Context) {
	test(this)
}
`
var f2 = `func Setup(debug bool) *Server {
	return &Server{}
}
`

var br1 = `sample bot response:
'''pkg/server/prompt_handler.go
%s'''
some other text:
'''pkg/server/server.go
%s'''
`

var _ = Describe("Parse BotResponse", func() {
	Context("with a well-formed bot response", func() {
		var parser preprocessors.Parser
		var response string

		BeforeEach(func() {
			response = fmt.Sprintf(br1, f1, f2)
			parser = preprocessors.Parser{CodeMap: make(preprocessors.SourceCodeMap)}
			Expect(parser.CodeMap).To(BeEmpty())
		})
		It("should match well-formed file paths", func() {
			Expect(preprocessors.IsValidFilePath("pkg/server/prompt_handler.go")).To(BeTrue())
			Expect(preprocessors.IsValidFilePath("pkg/server/server.go")).To(BeTrue())
			Expect(preprocessors.IsValidFilePath("/etc/config/cfg.yaml")).To(BeTrue())
			Expect(preprocessors.IsValidFilePath("C:\\Windows\\Sucks\\cfg.yaml")).To(BeFalse())
		})
		It("should successfully extract the correct content to the source code map", func() {
			Expect(parser.ParseBotResponse(response)).ShouldNot(HaveOccurred())
			Expect(len(parser.CodeMap)).To(Equal(2))
			Expect(parser.CodeMap["pkg/server/prompt_handler.go"]).To(Equal(f1))
			Expect(parser.CodeMap["pkg/server/server.go"]).To(Equal(f2))
		})
	})
	Context("when no valid code snippets are found", func() {
		It("should return an empty map", func() {
			parser := preprocessors.Parser{CodeMap: make(preprocessors.SourceCodeMap)}
			Expect(parser.ParseBotResponse("some text")).ShouldNot(HaveOccurred())
			Expect(parser.CodeMap).Should(BeEmpty())
		})
	})
	Context("with malformed code snippets", func() {
		It("should never match when file path is malformed", func() {
			parser := preprocessors.Parser{CodeMap: make(preprocessors.SourceCodeMap)}
			Expect(parser.ParseBotResponse("'''server\\prompt_handler.go\nsome text\n'''")).
				ShouldNot(HaveOccurred())
			Expect(parser.CodeMap).Should(BeEmpty())
		})
		It("should never match when file path is missing", func() {
			parser := preprocessors.Parser{CodeMap: make(preprocessors.SourceCodeMap)}
			Expect(parser.ParseBotResponse("'''\nsome text\nand more text.\n'''")).
				ShouldNot(HaveOccurred())
			Expect(parser.CodeMap).Should(BeEmpty())
		})
	})
})

var _ = Describe("Prompt Parser", func() {
	Context("with a well-formed user prompt", func() {
		var parser preprocessors.Parser
		BeforeEach(func() {
			parser = preprocessors.Parser{CodeMap: make(preprocessors.SourceCodeMap)}
		})

		It("should successfully fill in the correct content from the source code map", func() {
			prompt := fmt.Sprintf(br1, "", "")
			parser.ParsePrompt(prompt)
			Expect(parser.CodeMap).NotTo(BeEmpty())
			_, found := parser.CodeMap["pkg/server/prompt_handler.go"]
			Expect(found).To(BeTrue())
			_, found = parser.CodeMap["pkg/server/server.go"]
			Expect(found).To(BeTrue())
			parser.CodeMap["pkg/server/prompt_handler.go"] = f1
			parser.CodeMap["pkg/server/server.go"] = f2
			parsed, err := parser.FillPrompt(prompt)
			Expect(err).ShouldNot(HaveOccurred())
			expected := fmt.Sprintf(br1, f1, f2)
			Expect(parsed).To(Equal(expected))
		})
		It("should return the same prompt when no valid code snippets are found", func() {
			prompt := `some text
and more text:
'''
some_code: "that has no file path"
  should: "not be parsed"
'''
`
			parsed, err := parser.FillPrompt(prompt)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(parsed).To(Equal(prompt))
		})
	})
	Context("with malformed code snippets", func() {
		var parser preprocessors.Parser
		BeforeEach(func() {
			parser = preprocessors.Parser{CodeMap: make(preprocessors.SourceCodeMap)}
			parser.CodeMap["pkg/server/prompt_handler.go"] = f1
			parser.CodeMap["pkg/server/server.go"] = f2
			Expect(len(parser.CodeMap)).To(Equal(2))
		})
		PIt("should return an error when missing closing triple-quotes")
		PIt("should return an error when the file path is malformed")
		It("should not return an error when the code snippet is inserted manually", func() {
			prompt := `some text
and more text:
'''foo/bar.yaml
some_code: "that does not match a file path"
  should: "not fail"
'''
`
			parser.ParsePrompt(prompt)
			actual, err := parser.FillPrompt(prompt)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actual).To(Equal(prompt))
		})
		It("should return an error when file path does not exist", func() {
			prompt := `some text
This shouldn't work:
'''foo/bar.yaml
'''
`
			parser.ParsePrompt(prompt)
			// doing nothing here, is equivalent to not finding the file path
			_, err := parser.FillPrompt(prompt)
			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("FilesystemStore", func() {
	Context("Getting files from the filesystem", func() {
		var store preprocessors.CodeStoreHandler
		var codeMap preprocessors.SourceCodeMap
		BeforeEach(func() {
			codeMap = make(preprocessors.SourceCodeMap)
			// Copies the files from the testdata directory into a new temporary directory
			// and uses that as the working directory as the source for the tests.
			src, dest, err := SetupTestFiles()
			Expect(err).ShouldNot(HaveOccurred())
			store = preprocessors.NewFilesystemStore(src, dest)
			// Fills in the source code map with the file paths in the test directory,
			// but no contents
			Expect(FillCodemap(TestdataDir, codeMap, false)).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			// Removes the temporary directory
			Cleanup(store.(*preprocessors.FilesystemStore))
		})
		It("Should fill in the map with the correct file content", func() {
			Expect(store.GetSourceCode(&codeMap)).ToNot(HaveOccurred())
			Expect(len(codeMap)).To(Equal(7))
			for name, content := range codeMap {
				data, err := os.ReadFile(filepath.Join(TestdataDir, name))
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(string(data)))
			}
		})
		It("Should return error if file not found", func() {
			codeMap["foo/bar.yaml"] = ""
			Expect(store.GetSourceCode(&codeMap)).To(HaveOccurred())
		})
		It("Should correctly fill in files in subfolders", func() {
			data, _ := os.ReadFile(filepath.Join(TestdataDir, "misc/data/lvl/deep.txt"))
			Expect(store.GetSourceCode(&codeMap)).ToNot(HaveOccurred())
			Expect(codeMap["misc/data/lvl/deep.txt"]).To(Equal(string(data)))
		})
	})
	Context("Putting files to the filesystem", func() {
		var destDir string
		var err error
		var store preprocessors.CodeStoreHandler
		var codeMap = make(preprocessors.SourceCodeMap)

		//var store preprocessors.CodeStoreHandler
		BeforeEach(func() {
			// Creates a new temporary directory as the destination directory
			_, destDir, err = SetupTestFiles()
			Expect(err).ShouldNot(HaveOccurred())
			store = preprocessors.NewFilesystemStore("/foo/bar", destDir)
			// Fills in the source code map with test data
			Expect(FillCodemap(TestdataDir, codeMap, true)).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			// Removes the temporary directory
			Cleanup(store.(*preprocessors.FilesystemStore))
		})

		It("Should match file contents when content is saved via PutSourceCode", func() {
			err := store.PutSourceCode(codeMap)
			Expect(err).ToNot(HaveOccurred())
			isEqual, err := CompareDirs(TestdataDir, destDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(isEqual).To(BeTrue())
		})
		It("Should correctly handle nested directories in the file path", func() {
			codeMap["misc/data/lvl/level-2/test.txt"] = "some deeply-nested content"
			err := store.PutSourceCode(codeMap)
			Expect(err).ToNot(HaveOccurred())
			// this code checks that the file was created in the correct directory
			_, err = os.Stat(filepath.Join(destDir, "misc/data/lvl/level-2/test.txt"))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
