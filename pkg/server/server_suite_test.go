package server_test

import (
	"fmt"
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
	defer sourceFile.Close()

	var destFile *os.File
	destFile, err = os.CreateTemp("", "test_config.*.yaml")
	if err != nil {
		return
	}
	defer destFile.Close()
	dest = destFile.Name()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return
	}
	return
}

var TestConfigLocation string

var _ = BeforeSuite(func() {
	curDir, _ := os.Getwd()
	var prefix string
	if strings.HasSuffix(curDir, "server") {
		prefix = "../.."
	} else {
		prefix = ".."
	}
	// Set up the test environment
	TestConfigLocation = strings.Join([]string{prefix, "testdata/test_config.yaml"}, "/")
	fmt.Println(">>> BeforeSuite - TestConfigLocation: ", TestConfigLocation)
})
