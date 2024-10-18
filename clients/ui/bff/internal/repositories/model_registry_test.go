package repositories

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestFetchAllModelRegistry", func() {
	Context("with existing model registries", Ordered, func() {

		It("should retrieve the get all service successfully", func() {

			By("fetching all model registries in the repository")
			modelRegistryRepository := NewModelRegistryRepository()
			registries, err := modelRegistryRepository.FetchAllModelRegistries(k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model registries")
			expectedRegistries := []models.ModelRegistryModel{
				{Name: "model-registry", Description: "Model Registry Description", DisplayName: "Model Registry"},
				{Name: "model-registry-bella", Description: "Model Registry Bella description", DisplayName: "Model Registry Bella"},
				{Name: "model-registry-dora", Description: "Model Registry Dora description", DisplayName: "Model Registry Dora"},
			}
			Expect(registries).To(ConsistOf(expectedRegistries))
		})
	})
})
