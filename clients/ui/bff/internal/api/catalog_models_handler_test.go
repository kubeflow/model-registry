package api

import (
	"net/http"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAllCatalogModelsAcrossSourcesHandler", func() {
	Context("testing Get All Catalog Models Across Sources Handler", Ordered, func() {

		It("should retrieve all models for a source", func() {
			By("fetching all models for a source")
			data := mocks.GetCatalogModelListMock()
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			expected := CatalogModelListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[CatalogModelListEnvelope](http.MethodGet, "/api/v1/model_catalog/models?namespace=kubeflow&source=sample-source", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model sources")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(actual.Data.PageSize).To(Equal(expected.Data.PageSize))
			Expect(actual.Data.NextPageToken).To(Equal(expected.Data.NextPageToken))
			Expect(len(actual.Data.Items)).To(Equal(len(expected.Data.Items)))
			Expect(actual.Data.Items).To(Equal(expected.Data.Items))
		})

	})
})
