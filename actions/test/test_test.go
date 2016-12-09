package test_test

import (
	"bytes"
	"testing"

	"github.com/mikesimons/yaml-dsl/actions/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestAction Suite")
}

var _ = Describe("TestAction", func() {
	Describe("Prototype", func() {
		It("should return instance of TestAction", func() {
			subject := test.Prototype()
			testAction := (*test.TestAction)(nil)
			Expect(subject).Should(BeAssignableToTypeOf(testAction))
		})
	})

	Describe("Execute", func() {
		It("should print a dump of itself", func() {
			buf := &bytes.Buffer{}
			subject := &test.TestAction{}
			subject.Stdout = buf
			subject.Test = "This is a test"

			subject.Execute()

			Expect(buf.String()).Should(ContainSubstring("This is a test"))
		})

		It("should always return success", func() {
			buf := &bytes.Buffer{}
			subject := &test.TestAction{}
			subject.Stdout = buf

			result, err := subject.Execute()

			Expect(err).Should(BeNil())
			Expect(result.Success).Should(BeTrue())
		})
	})
})
