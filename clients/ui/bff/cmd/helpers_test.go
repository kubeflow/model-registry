package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("newOriginParser helper function", func() {
	var originParser func(s string) error
	var allowList []string

	BeforeEach(func() {
		allowList = []string{}
		originParser = newOriginParser(&allowList, "")
	})

	It("should parse a valid string list with 1 item", func() {
		expected := []string{"https://test.com"}

		err := originParser("https://test.com")

		Expect(err).NotTo(HaveOccurred())
		Expect(allowList).To(Equal(expected))
	})

	It("should parse a valid string list with 2 items", func() {
		expected := []string{"https://test.com", "https://test2.com"}

		err := originParser("https://test.com,https://test2.com")

		Expect(err).NotTo(HaveOccurred())
		Expect(allowList).To(Equal(expected))
	})

	It("should parse a valid string list with 2 items and extra spaces", func() {
		expected := []string{"https://test.com", "https://test2.com"}

		err := originParser("https://test.com,    https://test2.com")

		Expect(err).NotTo(HaveOccurred())
		Expect(allowList).To(Equal(expected))
	})

	It("should parse an empty string", func() {
		err := originParser("")

		Expect(err).NotTo(HaveOccurred())
		Expect(allowList).To(BeEmpty())
	})

	It("should parse the wildcard string", func() {
		expected := []string{"*"}

		err := originParser("*")

		Expect(err).NotTo(HaveOccurred())
		Expect(allowList).To(Equal(expected))
	})

})

func TestMainHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main helpers suite")
}
