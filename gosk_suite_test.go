package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGosk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gosk Suite")
}
