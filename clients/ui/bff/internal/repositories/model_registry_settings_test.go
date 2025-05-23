package repositories

import (
	"context"

	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestModelRegistrySettingsRepository", func() {
	Context("fetching groups from the stub", Ordered, func() {

		It("should return the mocked group list", func() {
			By("initializing the repository and client")
			repo := NewModelRegistrySettingsRepository()
			k8sClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			By("fetching groups")
			groups, err := repo.GetGroups(context.Background(), k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("verifying the returned group models")
			expected := []models.GroupModel{
				{Name: "dora-group-mock"},
				{Name: "bella-group-mock"},
			}
			Expect(groups).To(ConsistOf(expected))
		})
	})
})
