package inferenceservicecontroller_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	inferenceservicecontroller "github.com/kubeflow/model-registry/pkg/inferenceservice-controller"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	authv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	inferenceServiceIDLabel = "modelregistry.kubeflow.org/inference-service-id"
	registeredModelIDLabel  = "modelregistry.kubeflow.org/registered-model-id"
	modelVersionIDLabel     = "modelregistry.kubeflow.org/model-version-id"
	namespaceLabel          = "modelregistry.kubeflow.org/namespace"
	nameLabel               = "modelregistry.kubeflow.org/name"
	skipTLSVerify           = true
	urlAnnotation           = "modelregistry.kubeflow.org/url"
	finalizerLabel          = "modelregistry.kubeflow.org/finalizer"
	defaultNamespace        = "default"
	accessToken             = ""
	kserveVersion           = "v0.12.1"
	kserveCRDParamUrl       = "https://raw.githubusercontent.com/kserve/kserve/refs/tags/%s/config/crd/serving.kserve.io_inferenceservices.yaml"
	testCRDLocalPath        = "./testdata/crd"
)

var (
	cli          client.Client
	envTest      *envtest.Environment
	ctx          context.Context
	cancel       context.CancelFunc
	mrMockServer *httptest.Server
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "InferenceService Controller Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithCancel(context.Background())

	// Initialize logger
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.TimeEncoderOfLayout(time.RFC3339),
	}
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseFlagOptions(&opts)))

	// Download CRDs
	crdUrl := fmt.Sprintf(kserveCRDParamUrl, kserveVersion)

	crdPath := filepath.Join(testCRDLocalPath, "serving.kserve.io_inferenceservices.yaml")

	if err := os.MkdirAll(filepath.Dir(crdPath), 0755); err != nil {
		Fail(err.Error())
	}

	if err := DownloadFile(crdUrl, crdPath); err != nil {
		Fail(err.Error())
	}

	// Initialize test environment:
	By("Bootstrapping test environment")
	envTest = &envtest.Environment{
		CRDInstallOptions: envtest.CRDInstallOptions{
			Paths:              []string{filepath.Join("testdata", "crd")},
			ErrorIfPathMissing: true,
			CleanUpAfterUse:    false,
		},
	}

	cfg, err := envTest.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// Register API objects:
	testScheme := runtime.NewScheme()
	RegisterSchemes(testScheme)

	// Initialize Kubernetes client
	cli, err = client.New(cfg, client.Options{Scheme: testScheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(cli).NotTo(BeNil())

	// Setup controller manager
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:         testScheme,
		LeaderElection: false,
		Metrics: server.Options{
			BindAddress: "0",
		},
	})

	Expect(err).NotTo(HaveOccurred())

	const ModelRegistrySVCPath = "./testdata/deploy/model-registry-svc.yaml"

	mrSvc := &corev1.Service{}
	Expect(ConvertFileToStructuredResource(ModelRegistrySVCPath, mrSvc)).To(Succeed())

	mrSvc.SetNamespace(defaultNamespace)

	if err := cli.Create(ctx, mrSvc); err != nil && !errors.IsAlreadyExists(err) {
		Fail(err.Error())
	}

	inferenceServiceController := inferenceservicecontroller.NewInferenceServiceController(
		cli,
		ctrl.Log.WithName("controllers").WithName("ModelRegistry-InferenceService-Controller"),
		skipTLSVerify,
		accessToken,
		inferenceServiceIDLabel,
		registeredModelIDLabel,
		modelVersionIDLabel,
		namespaceLabel,
		nameLabel,
		urlAnnotation,
		finalizerLabel,
		defaultNamespace,
	)

	mrMockServer = ModelRegistryDefaultMockServer()

	inferenceServiceController.OverrideHTTPClient(&http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				if strings.Contains(req.URL.String(), "svc.cluster.local") {
					url, err := url.Parse(mrMockServer.URL)
					if err != nil {
						return nil, err
					}

					logf.Log.Info("Proxying request", "request", req.URL)

					proxyUrl, err := http.ProxyURL(url)(req)

					logf.Log.Info("Proxying request", "proxyUrl", proxyUrl)

					return proxyUrl, err
				}

				return req.URL, nil
			},
		},
	})

	err = inferenceServiceController.SetupWithManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	// Start the manager
	go func() {
		defer GinkgoRecover()
		err = mgr.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "Failed to run manager")
	}()
}, 60)

var _ = AfterSuite(func() {
	cancel()
	By("Tearing down the test environment")
	err := envTest.Stop()
	Expect(err).NotTo(HaveOccurred())
	By("Stopping the Model Registry mock server")
	mrMockServer.Close()

	// Clean up CRDs
	err = os.Remove(filepath.Join(testCRDLocalPath, "serving.kserve.io_inferenceservices.yaml"))
	Expect(err).NotTo(HaveOccurred())

})

func ModelRegistryDefaultMockServer() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/api/model_registry/v1alpha3/serving_environments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)

			res := `{
				"id": "1",
				"name": "default"
			}`

			_, err := w.Write([]byte(res))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	handler.HandleFunc("/api/model_registry/v1alpha3/inference_services/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res := `{
			"id": "1"
		}`

		_, err := w.Write([]byte(res))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	handler.HandleFunc("/api/model_registry/v1alpha3/inference_services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)

			res := `{
				"id": "1"
			}`

			_, err := w.Write([]byte(res))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewServer(handler)
}

func RegisterSchemes(s *runtime.Scheme) {
	utilruntime.Must(clientgoscheme.AddToScheme(s))
	utilruntime.Must(kservev1beta1.AddToScheme(s))
	utilruntime.Must(corev1.AddToScheme(s))
	utilruntime.Must(authv1.AddToScheme(s))
}

func ConvertFileToStructuredResource(path string, out runtime.Object) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return ConvertYamlToStructuredResource(data, out)
}

func ConvertYamlToStructuredResource(content []byte, out runtime.Object) error {
	s := runtime.NewScheme()

	RegisterSchemes(s)

	decoder := serializer.NewCodecFactory(s).UniversalDeserializer().Decode
	_, _, err := decoder(content, nil, out)

	return err
}

func DownloadFile(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
