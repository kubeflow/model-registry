package repositories

import (
	"context"
	"errors"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Model transfer jobs repository", func() {

	Describe("GetAllModelTransferJobs", func() {
		It("excludes jobs with DeletionTimestamp set", func() {
			realClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			now := metav1.Now()
			jobActive := batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "active-job",
					Namespace: "kubeflow",
					Labels: map[string]string{
						"modelregistry.kubeflow.org/job-type":            "async-upload",
						"modelregistry.kubeflow.org/model-registry-name": "model-registry",
						"modelregistry.kubeflow.org/job-id":              "id-active",
					},
					CreationTimestamp: metav1.Now(),
				},
				Status: batchv1.JobStatus{},
			}
			jobDeleting := batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "deleting-job",
					Namespace:         "kubeflow",
					DeletionTimestamp: &now,
					Labels: map[string]string{
						"modelregistry.kubeflow.org/job-type":            "async-upload",
						"modelregistry.kubeflow.org/model-registry-name": "model-registry",
						"modelregistry.kubeflow.org/job-id":              "id-deleting",
					},
					CreationTimestamp: metav1.Now(),
				},
				Status: batchv1.JobStatus{},
			}

			fakeClient := &jobListInjectorClient{
				KubernetesClientInterface: realClient,
				jobList: &batchv1.JobList{
					Items: []batchv1.Job{jobDeleting, jobActive},
				},
			}

			repo := NewModelRegistryRepository()
			list, err := repo.GetAllModelTransferJobs(
				mocks.NewMockSessionContextNoParent(),
				fakeClient,
				"kubeflow",
				"model-registry",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(list.Items).To(HaveLen(1))
			Expect(list.Items[0].Name).To(Equal("active-job"))
		})
	})

	Describe("registryOriginOnly", func() {
		It("returns scheme and host only (no port or path)", func() {
			Expect(registryOriginOnly("https://host.example.com:8080/api/path")).To(Equal("https://host.example.com"))
			Expect(registryOriginOnly("http://registry:9000/v2")).To(Equal("http://registry"))
		})
		It("defaults scheme to http when missing", func() {
			Expect(registryOriginOnly("//host.example.com/path")).To(Equal("http://host.example.com"))
		})
	})

	Describe("getModelRegistryAddress", func() {
		It("returns ErrModelRegistryNotFound when registry Get returns NotFound", func() {
			notFoundClient := &getServiceNotFoundClient{
				KubernetesClientInterface: nil,
			}
			repo := NewModelRegistryRepository()
			addr, err := repo.getModelRegistryAddress(
				mocks.NewMockSessionContextNoParent(),
				notFoundClient,
				"ns",
				"model-registry",
			)
			Expect(err).To(HaveOccurred())
			Expect(addr).To(BeEmpty())
			Expect(errors.Is(err, ErrModelRegistryNotFound)).To(BeTrue())
		})
	})

	Describe("buildModelMetadataConfigMap", func() {
		It("sets model_source_* keys and model_source_name to job name", func() {
			payload := models.ModelTransferJob{
				ModelVersionName: "v1",
				Author:           "author",
			}
			cm := buildModelMetadataConfigMap("cm-name", "ns", payload, "job-id", "my-job-name")
			Expect(cm.Data).NotTo(BeNil())
			Expect(cm.Data["ModelArtifact.model_source_kind"]).To(Equal("Job"))
			Expect(cm.Data["ModelArtifact.model_source_class"]).To(Equal("async-upload"))
			Expect(cm.Data["ModelArtifact.model_source_group"]).To(Equal("batch/v1"))
			Expect(cm.Data["ModelArtifact.model_source_name"]).To(Equal("my-job-name"))
		})
	})
})

// jobListInjectorClient returns a custom JobList from GetAllModelTransferJobs; other methods delegate.
type jobListInjectorClient struct {
	kubernetes.KubernetesClientInterface
	jobList *batchv1.JobList
}

func (c *jobListInjectorClient) GetAllModelTransferJobs(ctx context.Context, namespace, modelRegistryID string) (*batchv1.JobList, error) {
	return c.jobList, nil
}

// getServiceNotFoundClient returns NotFound from GetServiceDetailsByName; used to test getModelRegistryAddress.
type getServiceNotFoundClient struct {
	kubernetes.KubernetesClientInterface
}

func (c *getServiceNotFoundClient) GetServiceDetailsByName(ctx context.Context, namespace, serviceName, serviceType string) (kubernetes.ServiceDetails, error) {
	return kubernetes.ServiceDetails{}, apierrors.NewNotFound(schema.GroupResource{Resource: "services"}, serviceName)
}

// Stub the rest of the interface so the type compiles (getModelRegistryAddress only calls GetModelRegistry -> GetServiceDetailsByName).
// These must not be called in the getModelRegistryAddress test.
func (c *getServiceNotFoundClient) GetServiceNames(ctx context.Context, namespace string) ([]string, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetServiceDetails(ctx context.Context, namespace string) ([]kubernetes.ServiceDetails, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetNamespaces(ctx context.Context, identity *kubernetes.RequestIdentity) ([]corev1.Namespace, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) CanListServicesInNamespace(ctx context.Context, identity *kubernetes.RequestIdentity, namespace string) (bool, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) CanAccessServiceInNamespace(ctx context.Context, identity *kubernetes.RequestIdentity, namespace, serviceName string) (bool, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetSelfSubjectRulesReview(ctx context.Context, identity *kubernetes.RequestIdentity, namespace string) ([]string, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) IsClusterAdmin(identity *kubernetes.RequestIdentity) (bool, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) BearerToken() (string, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetUser(identity *kubernetes.RequestIdentity) (string, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetGroups(ctx context.Context) ([]string, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetAllCatalogSourceConfigs(ctx context.Context, namespace string) (corev1.ConfigMap, corev1.ConfigMap, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) UpdateCatalogSourceConfig(ctx context.Context, namespace string, configMap *corev1.ConfigMap) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) PatchSecret(ctx context.Context, namespace string, secretName string, data map[string]string) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) DeleteSecret(ctx context.Context, namespace string, secretName string) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetAllModelTransferJobs(ctx context.Context, namespace, modelRegistryID string) (*batchv1.JobList, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) CreateModelTransferJob(ctx context.Context, namespace string, job *batchv1.Job) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) DeleteModelTransferJob(ctx context.Context, namespace string, jobName string) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) CreateConfigMap(ctx context.Context, namespace string, configMap *corev1.ConfigMap) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) DeleteConfigMap(ctx context.Context, namespace string, name string) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetModelTransferJob(ctx context.Context, namespace string, jobName string) (*batchv1.Job, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetConfigMap(ctx context.Context, namespace string, name string) (*corev1.ConfigMap, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) GetSecret(ctx context.Context, namespace string, name string) (*corev1.Secret, error) {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) PatchSecretOwnerReference(ctx context.Context, namespace string, name string, ownerRef metav1.OwnerReference) error {
	panic("not used in getModelRegistryAddress test")
}
func (c *getServiceNotFoundClient) PatchConfigMapOwnerReference(ctx context.Context, namespace string, name string, ownerRef metav1.OwnerReference) error {
	panic("not used in getModelRegistryAddress test")
}
