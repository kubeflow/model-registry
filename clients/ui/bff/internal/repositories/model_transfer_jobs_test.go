package repositories

import (
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("resolveAsyncUploadImage", func() {
	var (
		client  k8s.KubernetesClientInterface
		mockCtx = mocks.NewMockSessionContextNoParent()
	)

	BeforeEach(func() {
		var err error
		client, err = kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when not in federated mode", func() {
		It("should return the default image", func() {
			repo := NewModelRegistryRepository(false, "")
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal(DefaultAsyncUploadImage))
		})
	})

	Context("when in federated mode with empty namespace", func() {
		It("should return the default image", func() {
			repo := NewModelRegistryRepository(true, "")
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal(DefaultAsyncUploadImage))
		})
	})

	Context("when in federated mode with ConfigMap missing", func() {
		It("should fall back to the default image", func() {
			repo := NewModelRegistryRepository(true, "bento-namespace")
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal(DefaultAsyncUploadImage))
		})
	})

	Context("when in federated mode with ConfigMap present", func() {
		const testNamespace = "kubeflow"

		AfterEach(func() {
			_ = client.DeleteConfigMap(mockCtx, testNamespace, asyncUploadConfigMapName)
		})

		It("should return the configured image when the key is set", func() {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      asyncUploadConfigMapName,
					Namespace: testNamespace,
				},
				Data: map[string]string{
					asyncUploadConfigMapKey: "registry.example.com/custom-image:v1",
				},
			}
			_, err := client.CreateConfigMap(mockCtx, testNamespace, cm)
			Expect(err).NotTo(HaveOccurred())

			repo := NewModelRegistryRepository(true, testNamespace)
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal("registry.example.com/custom-image:v1"))
		})

		It("should fall back to the default image when the key is missing", func() {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      asyncUploadConfigMapName,
					Namespace: testNamespace,
				},
				Data: map[string]string{
					"some-other-key": "some-value",
				},
			}
			_, err := client.CreateConfigMap(mockCtx, testNamespace, cm)
			Expect(err).NotTo(HaveOccurred())

			repo := NewModelRegistryRepository(true, testNamespace)
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal(DefaultAsyncUploadImage))
		})

		It("should fall back to the default image when the key is empty", func() {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      asyncUploadConfigMapName,
					Namespace: testNamespace,
				},
				Data: map[string]string{
					asyncUploadConfigMapKey: "",
				},
			}
			_, err := client.CreateConfigMap(mockCtx, testNamespace, cm)
			Expect(err).NotTo(HaveOccurred())

			repo := NewModelRegistryRepository(true, testNamespace)
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal(DefaultAsyncUploadImage))
		})

		It("should fall back to the default image when the key is whitespace-only", func() {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      asyncUploadConfigMapName,
					Namespace: testNamespace,
				},
				Data: map[string]string{
					asyncUploadConfigMapKey: "   ",
				},
			}
			_, err := client.CreateConfigMap(mockCtx, testNamespace, cm)
			Expect(err).NotTo(HaveOccurred())

			repo := NewModelRegistryRepository(true, testNamespace)
			img := repo.resolveAsyncUploadImage(mockCtx, client)
			Expect(img).To(Equal(DefaultAsyncUploadImage))
		})
	})
})
