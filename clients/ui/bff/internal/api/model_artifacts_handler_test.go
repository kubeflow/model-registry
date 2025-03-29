package api

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestArtifactsHandler", func() {
	Context("testing artifacts", Ordered, func() {
		Context("successful operations", func() {
			It("should retrieve an artifact", func() {
				By("fetching a specific artifact")
				_ = gofakeit.Seed(123)
				data := mocks.GenerateMockArtifact()
				expected := ArtifactEnvelope{Data: &data}

				_ = gofakeit.Seed(123)
				requestIdentity := kubernetes.RequestIdentity{
					UserID: "user@example.com",
				}
				actual, rs, err := setupApiTest[ArtifactEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/artifacts/1?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
				Expect(err).NotTo(HaveOccurred())

				By("should match the expected artifact")
				Expect(rs.StatusCode).To(Equal(http.StatusOK))
				Expect(actual.Data.ModelArtifact.GetName()).To(Equal(expected.Data.ModelArtifact.GetName()))
				Expect(actual.Data.ModelArtifact.GetArtifactType()).To(Equal(expected.Data.ModelArtifact.GetArtifactType()))
				Expect(actual.Data.ModelArtifact.GetDescription()).To(Equal(expected.Data.ModelArtifact.GetDescription()))
			})

			It("should list all artifacts", func() {
				By("fetching all artifacts")
				_ = gofakeit.Seed(123)
				requestIdentity := kubernetes.RequestIdentity{
					UserID: "user@example.com",
				}
				actual, rs, err := setupApiTest[ArtifactListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/artifacts?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
				Expect(err).NotTo(HaveOccurred())

				By("should return success status and valid data")
				Expect(rs.StatusCode).To(Equal(http.StatusOK))
				Expect(actual.Data).NotTo(BeNil())
				Expect(actual.Data.Items).NotTo(BeEmpty())

				By("should contain valid artifacts")
				for _, item := range actual.Data.Items {
					Expect(*item.ModelArtifact.Name).NotTo(BeEmpty())
					Expect(*item.ModelArtifact.ArtifactType).NotTo(BeEmpty())
				}
			})
		})
	})
})
