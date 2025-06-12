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
			Expect(groups).To(HaveLen(2))

			// Verify first group
			firstGroup := groups[0]
			Expect(firstGroup.Metadata.Name).To(Equal("dora-group-mock"))
			Expect(firstGroup.Users).To(ConsistOf("dora-user@example.com", "dora-admin@example.com"))
			Expect(*firstGroup.APIVersion).To(Equal("user.openshift.io/v1"))
			Expect(*firstGroup.Kind).To(Equal("Group"))
			Expect(firstGroup).To(BeAssignableToTypeOf(models.Group{}))

			// Verify second group
			secondGroup := groups[1]
			Expect(secondGroup.Metadata.Name).To(Equal("bella-group-mock"))
			Expect(secondGroup.Users).To(ConsistOf("bella-user@example.com", "bella-maintainer@example.com"))
			Expect(*secondGroup.APIVersion).To(Equal("user.openshift.io/v1"))
			Expect(*secondGroup.Kind).To(Equal("Group"))
		})
	})
})
