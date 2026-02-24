package tests

import (
	"context"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CanNamespaceAccessRegistry (shared)", func() {
	Context("when default SA in job namespace has access to the registry", func() {
		It("returns true for dora-namespace default SA and model-registry-dora in dora-namespace", func() {
			ctx := context.Background()
			allowed, err := kubernetes.CanNamespaceAccessRegistry(ctx, clientset, logger, "dora-namespace", "model-registry-dora", "dora-namespace")
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue())
		})

		It("returns true for bella-namespace default SA and model-registry-bella in bella-namespace", func() {
			ctx := context.Background()
			allowed, err := kubernetes.CanNamespaceAccessRegistry(ctx, clientset, logger, "bella-namespace", "model-registry-bella", "bella-namespace")
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue())
		})

		It("returns true for kubeflow default SA and model-registry in kubeflow", func() {
			ctx := context.Background()
			allowed, err := kubernetes.CanNamespaceAccessRegistry(ctx, clientset, logger, "kubeflow", "model-registry", "kubeflow")
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue())
		})
	})

	Context("when default SA in job namespace has no access to the registry", func() {
		It("returns false when bella-namespace default SA checks access to model-registry-dora in dora-namespace", func() {
			ctx := context.Background()
			allowed, err := kubernetes.CanNamespaceAccessRegistry(ctx, clientset, logger, "bella-namespace", "model-registry-dora", "dora-namespace")
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeFalse())
		})

		It("returns false when job namespace default SA checks non-existent registry", func() {
			ctx := context.Background()
			allowed, err := kubernetes.CanNamespaceAccessRegistry(ctx, clientset, logger, "dora-namespace", "nonexistent-registry", "dora-namespace")
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeFalse())
		})
	})
})
