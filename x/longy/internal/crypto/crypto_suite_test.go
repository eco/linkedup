package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crypto Test Suite")
}

// BeforeTestRun sets up common data needed by all handler tests.
func BeforeTestRun() {
}
