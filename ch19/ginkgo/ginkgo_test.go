package ginkgo

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

var _ = Describe("Books", func() {
	var (
		longBook  string
		shortBook string
		pathches  *gomonkey.Patches
		ctl       *gomock.Controller
	)
	BeforeEach(func() {
		longBook = "long"
		shortBook = "short"

		ctl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		longBook = ""
		shortBook = ""

		ctl.Finish()
		pathches.Reset()
	})

	Describe("Add Books", func() {
		It("should return a list of books", func() {
			// Expect(true).To(Equal(true))

			Expect(longBook).To(Equal("long"))
		})

		It("should return a list of books", func() {
			// Expect(true).To(Equal(true))

			Expect(shortBook).To(Equal("short"))
		})
	})

	Describe("Delete Books", func() {

	})
})
