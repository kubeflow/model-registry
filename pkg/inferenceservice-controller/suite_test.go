package inferenceservicecontroller_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	inferenceservicecontroller "github.com/kubeflow/model-registry/pkg/inferenceservice-controller"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	authv1 "k8s.io/api/rbac/v1"
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
	serviceURLAnnotation    = "routing.kubeflow.org/external-address-rest"
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
		serviceURLAnnotation,
		defaultNamespace,
	)

	mrMockServer = ModelRegistryDefaultMockServer()

	mockUrl, _ := url.Parse(mrMockServer.URL)

	mrMockServer.Client().Transport = RewriteTransport{
		Transport: mrMockServer.Client().Transport,
		URL:       mockUrl,
	}

	inferenceServiceController.OverrideHTTPClient(mrMockServer.Client())

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

	servingEnvironments := make(map[string]*openapi.ServingEnvironment)
	inferenceServices := make(map[string]*openapi.InferenceService)

	handler.HandleFunc("/api/model_registry/v1alpha3/serving_environments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if r.Method == http.MethodPost {
			id := "1"

			senv := &openapi.ServingEnvironment{}

			senvJson, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if err := json.Unmarshal(senvJson, &senv); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			senv.Id = &id

			servingEnvironments[id] = senv

			w.WriteHeader(http.StatusCreated)

			res, err := json.Marshal(senv)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = w.Write([]byte(res))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	handler.HandleFunc("/api/model_registry/v1alpha3/inference_services/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := "1"

		if r.Method == http.MethodGet {
			isvc, ok := inferenceServices[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			res, err := json.Marshal(isvc)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = w.Write([]byte(res))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}

		if r.Method == http.MethodPatch {
			isvc := &openapi.InferenceService{}

			isvcvJson, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if err := json.Unmarshal(isvcvJson, &isvc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			isvc.Id = &id

			inferenceServices[id] = isvc

			w.WriteHeader(http.StatusCreated)

			res, err := json.Marshal(isvc)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = w.Write([]byte(res))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	})

	handler.HandleFunc("/api/model_registry/v1alpha3/inference_services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if r.Method == http.MethodPost {
			id := "1"

			isvc := &openapi.InferenceService{}

			isvcvJson, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if err := json.Unmarshal(isvcvJson, &isvc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			isvc.Id = &id

			inferenceServices[id] = isvc

			w.WriteHeader(http.StatusCreated)

			res, err := json.Marshal(isvc)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = w.Write([]byte(res))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewTLSServer(handler)
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

	//nolint:errcheck
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

type RewriteTransport struct {
	Transport http.RoundTripper
	URL       *url.URL
}

func (t RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.URL.Scheme
	req.URL.Host = t.URL.Host
	req.URL.Path = path.Join(t.URL.Path, req.URL.Path)
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(req)
}
