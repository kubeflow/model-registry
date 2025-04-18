package tests

import (
	"context"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	k8mocks "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes/k8mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("KubernetesClient SSAR Test", func() {
	Context("CanListServicesInNamespace", func() {

		It("should allow allowed user to access services", func() {
			By("should allow allowed user to access services")

			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[0].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			allowed, err := tokenK8client.CanListServicesInNamespace(ctx, identity, "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SSAR for Kubeflow User ID\"")
			Expect(allowed).To(BeTrue(), "Expected Kubeflow User ID to have access")
		})

		It("should check dora's access to namespaces", func() {
			By("performing SAR for dora user")
			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[1].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			allowed, err := tokenK8client.CanListServicesInNamespace(ctx, identity, "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SAR for unauthorized-dora@example.com")
			Expect(allowed).To(BeFalse(), "Expected doraNonAdmin@example.com to be denied access")

			allowed, err = tokenK8client.CanListServicesInNamespace(ctx, identity, "dora-namespace")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SSAR for unauthorized-dora@example.com")
			Expect(allowed).To(BeTrue(), "Expected doraNonAdmin@example.com ID to have access")
		})

		It("should deny access for another user", func() {
			By("performing SAR for another user")
			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[2].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			allowed, err := tokenK8client.CanListServicesInNamespace(ctx, identity, "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SSAR for unauthorized-dora@example.com")
			Expect(allowed).To(BeFalse(), "Expected unauthorized-dora@example.com to be denied access")
		})

	})

	Context("CanAccessServiceInNamespace", func() {

		It("should allow allowed user to access service", func() {
			By("should allow allowed user to access service")

			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[0].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			allowed, err := tokenK8client.CanAccessServiceInNamespace(ctx, identity, "kubeflow", "model-registry")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SSAR for Kubeflow User ID\"")
			Expect(allowed).To(BeTrue(), "Expected Kubeflow User ID to have access")
		})

		It("should check dora's access to services", func() {
			By("performing SAR for dora user")
			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[1].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			allowed, err := tokenK8client.CanAccessServiceInNamespace(ctx, identity, "dora-namespace", "model-registry-dora")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SSAR")
			Expect(allowed).To(BeTrue())

			allowed, err = tokenK8client.CanAccessServiceInNamespace(ctx, identity, "bella-namespace", "model-registry-bella")
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeFalse())
		})

	})

	Context("GetNamespaces", func() {

		It("should allow allowed user to get namespaces", func() {
			By("cluster admin should be allowed to get namespaces")

			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[0].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			ns, err := tokenK8client.GetNamespaces(ctx, identity)
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SSAR for Kubeflow User ID\"")
			Expect(ns).NotTo(BeEmpty(), "Expected Kubeflow User ID to have access")
		})

		It("should allow allowed user to get namespaces", func() {
			By("other users should not be allowed to get namespaces")

			identity := &kubernetes.RequestIdentity{
				Token: k8mocks.DefaultTestUsers[1].Token,
			}
			ctx := context.WithValue(context.Background(), constants.RequestIdentityKey, identity)

			kubernetesMockedTokenClientFactory, err := k8mocks.NewTokenClientFactory(clientset, restConfig, logger)
			Expect(err).NotTo(HaveOccurred())
			tokenK8client, err := kubernetesMockedTokenClientFactory.GetClient(ctx)
			Expect(err).NotTo(HaveOccurred())

			ns, err := tokenK8client.GetNamespaces(ctx, identity)
			Expect(ns).To(BeEmpty())
			Expect(err).To(HaveOccurred())
		})
	})
})
