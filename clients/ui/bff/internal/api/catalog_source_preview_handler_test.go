package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateCatalogSourcePreviewHandler", func() {
	Context("testing create catalog source preview handler", Ordered, func() {
		It("should create source preview for a source", func() {
			By("creating a source preview")
			data := mocks.CreateCatalogSourcePreviewMock()
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			requestBody := struct {
				Data models.CatalogSourcePreviewRequest `json:"data"`
			}{
				Data: models.CatalogSourcePreviewRequest{
					Type:           "yaml",
					IncludedModels: []string{},
					ExcludedModels: []string{},
					Properties: map[string]interface{}{
						"yaml": "models:\n  - name: test-model",
					},
				},
			}
			bodyBytes, _ := json.Marshal(requestBody)

			expected := CatalogSourcePreviewEnvelope{Data: &data}
			actual, rs, err := setupApiTest[CatalogSourcePreviewEnvelope](http.MethodPost, "/api/v1/settings/model_catalog/source_preview?namespace=kubeflow", bytes.NewReader(bodyBytes), kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected source preview")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(actual.Data.PageSize).To(Equal(expected.Data.PageSize))
			Expect(actual.Data.NextPageToken).To(Equal(expected.Data.NextPageToken))
			Expect(actual.Data.Items).To(Equal(expected.Data.Items))

		})
	})
})
