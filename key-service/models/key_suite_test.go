package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestStoredKey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StoredKey Suite")
}
