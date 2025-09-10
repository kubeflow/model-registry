package api

import (
	"net/http"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestGetAllCatalogSourcesHandler", func() {
	Context("testing Catalog Sources Handler", Ordered, func() {

		It("should retrieve all catalog sources", func() {
			By("fetching all catalog sources")
			data := mocks.GetCatalogSourceListMock()
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			expected := CatalogSourceListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[CatalogSourceListEnvelope](http.MethodGet, "/api/v1/model_catalog/sources?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected catalog sources")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(actual.Data.PageSize).To(Equal(expected.Data.PageSize))
			Expect(actual.Data.NextPageToken).To(Equal(expected.Data.NextPageToken))
			Expect(len(actual.Data.Items)).To(Equal(len(expected.Data.Items)))
			Expect(actual.Data.Items).To(Equal(expected.Data.Items))
		})

		It("should retrieve all catalog sources filtered by name", func() {
			By("fetching catalog sources filtered by name")
			data := mocks.GetCatalogSourceListMock()
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			actual, rs, err := setupApiTest[CatalogSourceListEnvelope](http.MethodGet, "/api/v1/model_catalog/sources?namespace=kubeflow&name=dora", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected catalog sources")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(len(actual.Data.Items)).To(Equal(1))
			//dora is the [4]item in the mock data
			Expect(actual.Data.Items[0].Name).To(Equal(data.Items[4].Name))
		})

		It("should retrieve catalog sources model", func() {
			By("fetching catalog sources models")
			data := mocks.GetCatalogModelMocks()[0]
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			actual, rs, err := setupApiTest[CatalogModelEnvelope](http.MethodGet, "/api/v1/model_catalog/sources/source/models/model-name?namespace=kubeflow&name=dora", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			//our mocked version returns always the first model for any path. We just check here that handlers are working
			Expect(actual.Data.Name).To(Equal(data.Name))
		})

	})
})
