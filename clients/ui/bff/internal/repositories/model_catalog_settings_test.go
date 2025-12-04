package repositories

import (
	"context"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModelCatalogSettingRepository", func() {
	var (
		repo      *ModelCatalogSettingsRepository
		k8sClient k8s.KubernetesClientInterface
		ctx       context.Context
	)

	BeforeEach(func() {
		repo = NewModelCatalogSettingsRepository()
		var err error
		k8sClient, err = kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
		Expect(err).NotTo(HaveOccurred())
		ctx = mocks.NewMockSessionContextNoParent()
	})

	Describe("GetAllCatalogSourceConfigs", func() {
		It("should return merged source config from both default and user managed configMap", func() {
			catalogs, err := repo.GetAllCatalogSourceConfigs(ctx, k8sClient, "kubeflow")
			Expect(err).NotTo(HaveOccurred())
			Expect(catalogs).NotTo(BeEmpty())
		})

		It("deafult catalog should have isDefault=true", func() {
			catalogs, err := repo.GetAllCatalogSourceConfigs(ctx, k8sClient, "kubeflow")
			Expect(err).NotTo(HaveOccurred())
			var defaultCatalog *models.CatalogSourceConfig
			for _, c := range catalogs.Catalogs {
				if c.Id == "dora_ai_models" {
					defaultCatalog = &c
					break
				}
			}
			Expect(defaultCatalog).NotTo(BeNil())
			Expect(*defaultCatalog.IsDefault).To(BeTrue())
		})

		It("should merge, if user overrides the default source", func() {
			catalogs, err := repo.GetAllCatalogSourceConfigs(ctx, k8sClient, "kubeflow")
			Expect(err).NotTo(HaveOccurred())
			var mergedCatalog *models.CatalogSourceConfig
			for _, c := range catalogs.Catalogs {
				if c.Id == "dora_ai_models" {
					mergedCatalog = &c
					break
				}
			}

			Expect(mergedCatalog).NotTo(BeNil())
		})

	})
})
