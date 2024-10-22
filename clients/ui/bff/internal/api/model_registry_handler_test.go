package api

import (
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("TestModelRegistryHandler", func() {
	Context("fetching model registries", Ordered, func() {
		It("should retrieve the model registries successfully", func() {

			By("creating the test app")
			testApp := App{
				kubernetesClient: k8sClient,
				repositories:     repositories.NewRepositories(mockMRClient),
				logger:           logger,
			}

			By("creating the http test infrastructure")
			req, err := http.NewRequest(http.MethodGet, ModelRegistryListPath, nil)
			Expect(err).NotTo(HaveOccurred())
			rr := httptest.NewRecorder()

			By("creating the http request for the handler")
			testApp.ModelRegistryHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the model registries")
			var actual ModelRegistryListEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("should match the expected model registries")
			var expected = []models.ModelRegistryModel{
				{Name: "model-registry", Description: "Model Registry Description", DisplayName: "Model Registry"},
				{Name: "model-registry-bella", Description: "Model Registry Bella description", DisplayName: "Model Registry Bella"},
				{Name: "model-registry-dora", Description: "Model Registry Dora description", DisplayName: "Model Registry Dora"},
			}
			Expect(actual.Data).To(ConsistOf(expected))
		})

	})
})
