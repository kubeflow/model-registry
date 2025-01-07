package repositories

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestFetchAllModelRegistry", func() {
	Context("with existing model registries", Ordered, func() {

		It("should retrieve the get all kubeflow service successfully", func() {

			By("fetching all model registries in the repository")
			modelRegistryRepository := NewModelRegistryRepository()
			registries, err := modelRegistryRepository.GetAllModelRegistries(k8sClient, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model registries")
			expectedRegistries := []models.ModelRegistryModel{
				{Name: "model-registry", Description: "Model Registry Description", DisplayName: "Model Registry"},
				{Name: "model-registry-one", Description: "Model Registry One description", DisplayName: "Model Registry One"},
			}
			Expect(registries).To(ConsistOf(expectedRegistries))
		})

		It("should retrieve the get all dora-namespace service successfully", func() {

			By("fetching all model registries in the repository")
			modelRegistryRepository := NewModelRegistryRepository()
			registries, err := modelRegistryRepository.GetAllModelRegistries(k8sClient, "dora-namespace")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model registries")
			expectedRegistries := []models.ModelRegistryModel{
				{Name: "model-registry-dora", Description: "Model Registry Dora description", DisplayName: "Model Registry Dora"},
			}
			Expect(registries).To(ConsistOf(expectedRegistries))
		})

		It("should not retrieve namespaces", func() {

			By("fetching all model registries in the repository")
			modelRegistryRepository := NewModelRegistryRepository()
			registries, err := modelRegistryRepository.GetAllModelRegistries(k8sClient, "no-namespace")
			Expect(err).NotTo(HaveOccurred())

			By("should be empty")
			Expect(registries).To(BeEmpty())
		})
	})
})
