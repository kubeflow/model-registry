package tests

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenClientFactory ExtractRequestIdentity", func() {

	var factory *kubernetes.TokenClientFactory
	var header http.Header

	BeforeEach(func() {
		header = http.Header{}
	})

	Context("with Bearer prefix", func() {
		BeforeEach(func() {
			factory = &kubernetes.TokenClientFactory{
				Logger: nil,
				Header: config.DefaultAuthTokenHeader,
				Prefix: config.DefaultAuthTokenPrefix,
			}
		})

		It("should extract the token successfully", func() {
			header.Set("Authorization", "Bearer doratoken")

			identity, err := factory.ExtractRequestIdentity(header)
			Expect(err).NotTo(HaveOccurred())
			Expect(identity.Token).To(Equal("doratoken"))
		})

		It("should fail if prefix does not match", func() {
			header.Set("Authorization", "Token bellatoken")

			_, err := factory.ExtractRequestIdentity(header)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected token Header Authorization to start with Prefix \"Bearer \""))
		})

		It("should fail if prefix is missing", func() {
			header.Set("Authorization", "doratoken")

			_, err := factory.ExtractRequestIdentity(header)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected token Header Authorization to start with Prefix \"Bearer \""))
		})

		It("should fail if header is missing", func() {
			_, err := factory.ExtractRequestIdentity(header)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required Header: Authorization"))
		})
	})

	Context("with no prefix", func() {
		BeforeEach(func() {
			factory = &kubernetes.TokenClientFactory{
				Logger: nil,
				Header: "X-Forwarded-Access-Token",
				Prefix: "",
			}
		})

		It("should extract the raw token", func() {
			header.Set("X-Forwarded-Access-Token", "bellatoken")

			identity, err := factory.ExtractRequestIdentity(header)
			Expect(err).NotTo(HaveOccurred())
			Expect(identity.Token).To(Equal("bellatoken"))
		})

		It("should fail if header is missing", func() {
			_, err := factory.ExtractRequestIdentity(header)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required Header: X-Forwarded-Access-Token"))
		})

		It("should fail if header mismatch", func() {
			header.Set("X-WRONG-Access-Token", "bellatoken")

			_, err := factory.ExtractRequestIdentity(header)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required Header: X-Forwarded-Access-Token"))
		})
	})
})
