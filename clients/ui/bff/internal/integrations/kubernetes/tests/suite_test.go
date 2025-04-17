package tests

import (
	"context"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes/k8mocks"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	kubernetesMockedStaticClientFactory k8s.KubernetesClientFactory

	ctx        context.Context
	cancel     context.CancelFunc
	logger     *slog.Logger
	testEnv    *envtest.Environment
	clientset  kubernetes.Interface
	restConfig *rest.Config
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "API Suite")
}

var _ = BeforeSuite(func() {
	defer GinkgoRecover()
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.Background())

	By("bootstrapping test environment")
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	By("bootstrapping envtest")
	var err error

	testEnv, clientset, err = k8mocks.SetupEnvTest(k8mocks.TestEnvInput{
		Logger: logger,
		Ctx:    ctx,
		Cancel: cancel,
	})
	Expect(err).NotTo(HaveOccurred())
	restConfig = testEnv.Config

	By("creating factory mock client using shared envtest")
	kubernetesMockedStaticClientFactory, err = k8mocks.NewStaticClientFactory(clientset, logger)
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	By("shutting down the test environment")
	defer cancel()
	logger.Info("Stopping envtest control plane")
	if err := testEnv.Stop(); err != nil {
		logger.Error("failed to stop envtest", "error", err)
		Fail("Failed to stop envtest: " + err.Error())
	} else {
		logger.Info("envtest stopped successfully")
	}

})
