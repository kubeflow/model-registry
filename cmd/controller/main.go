package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strconv"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"github.com/kubeflow/model-registry/internal/controller/controllers"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	infrctrl "github.com/kubeflow/model-registry/pkg/inferenceservice-controller"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	var tlsOpts []func(*tls.Config)
	flag.StringVar(&metricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})

	// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
	// More info:
	// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/metrics/server
	// - https://book.kubebuilder.io/reference/metrics.html
	metricsServerOptions := metricsserver.Options{
		BindAddress:   metricsAddr,
		SecureServing: secureMetrics,
		TLSOpts:       tlsOpts,
	}

	if secureMetrics {
		// FilterProvider is used to protect the metrics endpoint with authn/authz.
		// These configurations ensure that only authorized users and service accounts
		// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
		// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/metrics/filters#WithAuthenticationAndAuthorization
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization

		// TODO(user): If CertDir, CertName, and KeyName are not specified, controller-runtime will automatically
		// generate self-signed certificates for the metrics server. While convenient for development and testing,
		// this setup is not recommended for production.
	}

	utilruntime.Must(kservev1beta1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsServerOptions,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a7d60e25.kubeflow.org",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if os.Getenv("INFERENCE_SERVICE_CONTROLLER") == "managed" {
		inferenceServiceController, err := setupInferenceServiceController(
			context.Background(),
			mgr,
			ctrl.GetConfigOrDie(),
		)
		if err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "InferenceService")
			os.Exit(1)
		}

		if err = (&controllers.InferenceServiceReconciler{
			Client:                     mgr.GetClient(),
			Scheme:                     mgr.GetScheme(),
			InferenceServiceController: inferenceServiceController,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "InferenceService")
			os.Exit(1)
		}
		// +kubebuilder:scaffold:builder
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupInferenceServiceController(ctx context.Context, mgr manager.Manager, cfg *rest.Config) (*infrctrl.InferenceServiceController, error) {
	namespaceLabel, err := getEnvOrFail("NAMESPACE_LABEL")
	if err != nil {
		return nil, err
	}

	nameLabel, err := getEnvOrFail("NAME_LABEL")
	if err != nil {
		return nil, err
	}

	urlAnnotation, err := getEnvOrFail("URL_ANNOTATION")
	if err != nil {
		return nil, err
	}

	inferenceServiceIDLabel, err := getEnvOrFail("INFERENCE_SERVICE_ID_LABEL")
	if err != nil {
		return nil, err
	}

	modelVersionIDLabel, err := getEnvOrFail("MODEL_VERSION_ID_LABEL")
	if err != nil {
		return nil, err
	}

	registeredModelIdLabel, err := getEnvOrFail("REGISTERED_MODEL_ID_LABEL")
	if err != nil {
		return nil, err
	}

	finalizer, err := getEnvOrFail("FINALIZER")
	if err != nil {
		return nil, err
	}

	serviceAnnotation, err := getEnvOrFail("SERVICE_ANNOTATION")
	if err != nil {
		return nil, err
	}

	registriesNamespace, err := getEnvOrFail("REGISTRIES_NAMESPACE")
	if err != nil {
		return nil, err
	}

	skipTLSVerify := getEnvAsBool("SKIP_TLS_VERIFY", false)

	return infrctrl.NewInferenceServiceController(
		mgr.GetClient(),
		log.FromContext(ctx).WithName("controllers").WithName("ModelRegistryInferenceService"),
		skipTLSVerify,
		cfg.BearerToken,
		inferenceServiceIDLabel,
		registeredModelIdLabel,
		modelVersionIDLabel,
		namespaceLabel,
		nameLabel,
		urlAnnotation,
		finalizer,
		serviceAnnotation,
		registriesNamespace,
	), nil
}

func getEnvOrFail(name string) (string, error) {
	valStr := os.Getenv(name)

	if valStr == "" {
		return "", fmt.Errorf("environment variable %s is required", name)
	}

	return valStr, nil
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := os.Getenv(name)

	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}
