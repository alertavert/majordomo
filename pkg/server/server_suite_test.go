package server_test

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

// MkTempConfigFile creates a temporary copy of the test config file, and returns its path.
func MkTempConfigFile(src string) (dest string, err error) {
	var sourceFile *os.File
	sourceFile, err = os.Open(src)
	if err != nil {
		return
	}
	defer func(sourceFile *os.File) {
		Ω(sourceFile.Close()).Should(Succeed())
	}(sourceFile)

	var destFile *os.File
	destFile, err = os.CreateTemp("", "test_config.*.yaml")
	if err != nil {
		return
	}
	defer func(destFile *os.File) {
		Ω(destFile.Close()).Should(Succeed())
	}(destFile)
	dest = destFile.Name()

	_, err = io.Copy(destFile, sourceFile)
	Ω(err).ShouldNot(HaveOccurred())
	return
}

var TestConfigLocation string

var _ = BeforeSuite(func() {
	// Silence the logs
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// Determine the prefix
	curDir, _ := os.Getwd()
	var prefix string
	if strings.HasSuffix(curDir, "server") {
		prefix = "../.."
	} else {
		prefix = ".."
	}
	// Set up the test environment
	TestConfigLocation = strings.Join([]string{prefix, "testdata/test_config.yaml"}, "/")
})
