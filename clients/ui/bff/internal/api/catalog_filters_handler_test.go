package api

import (
	"net/http"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestGetCatalogFilterListHandler", func() {
	Context("testing catalog filter list handler", Ordered, func() {

		It("should retrive filter option list", func() {
			By("fetching all filter option list")
			data := mocks.GetFilterOptionsListMock()
			requestIdentify := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			expected := CatalogFilterOptionsListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[CatalogFilterOptionsListEnvelope](http.MethodGet, "/api/v1/model_catalog/models/filter_options?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentify, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected catalog filter options")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data).To(Equal(expected.Data))
		})
	})
})
