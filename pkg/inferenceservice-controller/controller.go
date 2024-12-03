package inferenceservicecontroller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/kubeflow/model-registry/pkg/openapi"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type InferenceServiceController struct {
	client                        client.Client
	httpClient                    *http.Client
	log                           logr.Logger
	bearerToken                   string
	inferenceServiceIDLabel       string
	registeredModelIDLabel        string
	modelVersionIDLabel           string
	modelRegistryNamespaceLabel   string
	modelRegistryNameLabel        string
	modelRegistryURLAnnotation    string
	modelRegistryFinalizer        string
	defaultModelRegistryNamespace string
}

func NewInferenceServiceController(
	client client.Client,
	log logr.Logger,
	skipTLSVerify bool,
	bearerToken,
	isIDLabel,
	regModelIDLabel,
	modelVerIDLabel,
	mrNamespaceLabel,
	mrNameLabel,
	mrURLAnnotation,
	mrFinalizer,
	defaultMRNamespace string,
) *InferenceServiceController {
	httpClient := http.DefaultClient

	if skipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &InferenceServiceController{
		client:                        client,
		httpClient:                    httpClient,
		log:                           log,
		bearerToken:                   bearerToken,
		inferenceServiceIDLabel:       isIDLabel,
		registeredModelIDLabel:        regModelIDLabel,
		modelVersionIDLabel:           modelVerIDLabel,
		modelRegistryNamespaceLabel:   mrNamespaceLabel,
		modelRegistryNameLabel:        mrNameLabel,
		modelRegistryURLAnnotation:    mrURLAnnotation,
		modelRegistryFinalizer:        mrFinalizer,
		defaultModelRegistryNamespace: defaultMRNamespace,
	}
}

func (r *InferenceServiceController) OverrideHTTPClient(client *http.Client) {
	r.httpClient = client
}

// Reconcile performs the reconciliation of the model registry based on Kubeflow InferenceService CRs
func (r *InferenceServiceController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mrNamespace := r.defaultModelRegistryNamespace
	mrIs := &openapi.InferenceService{}
	mrApiCtx := context.Background()

	if r.bearerToken != "" {
		mrApiCtx = context.WithValue(context.Background(), openapi.ContextAccessToken, r.bearerToken)
	}

	// Initialize logger format
	log := r.log.WithValues("InferenceService", req.Name, "namespace", req.Namespace)

	isvc := &kservev1beta1.InferenceService{}
	err := r.client.Get(ctx, req.NamespacedName, isvc)
	if err != nil {
		if apierrs.IsNotFound(err) {
			log.V(1).Info("Stop ModelRegistry InferenceService reconciliation, ISVC not found.")
			return ctrl.Result{}, nil
		}

		log.Error(err, "Unable to fetch the InferenceService")
		return ctrl.Result{}, err
	}

	mrIsvcId, okMrIsvcId := isvc.Labels[r.inferenceServiceIDLabel]
	registeredModelId, okRegisteredModelId := isvc.Labels[r.registeredModelIDLabel]
	modelVersionId := isvc.Labels[r.modelVersionIDLabel]
	mrName, okMrName := isvc.Labels[r.modelRegistryNameLabel]
	mrUrl, okMrUrl := isvc.Annotations[r.modelRegistryURLAnnotation]

	if !okMrIsvcId && !okRegisteredModelId {
		// Early check: no model registry specific labels set in the ISVC, ignore the CR
		log.Error(fmt.Errorf("missing model registry specific label, unable to link ISVC to Model Registry, skipping InferenceService"), "Stop ModelRegistry InferenceService reconciliation")
		return ctrl.Result{}, nil
	}

	if !okMrName && !okMrUrl {
		// Early check: it's required to have the model registry name or url set in the ISVC
		log.Error(fmt.Errorf("missing model registry name or url, unable to link ISVC to Model Registry, skipping InferenceService"), "Stop ModelRegistry InferenceService reconciliation")
		return ctrl.Result{}, nil
	}

	if mrNSFromISVC, ok := isvc.Labels[r.modelRegistryNamespaceLabel]; ok {
		mrNamespace = mrNSFromISVC
	}

	log.Info("Creating model registry service..")
	mrApi, err := r.initModelRegistryService(ctx, log, mrName, mrNamespace, mrUrl)
	if err != nil {
		log.Error(err, "Unable to initialize Model Registry Service")
		return ctrl.Result{}, err
	}

	// Check if the InferenceService instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isMarkedToBeDeleted := isvc.GetDeletionTimestamp() != nil

	// Let's add a finalizer. Then, we can define some operations which should
	// occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !isMarkedToBeDeleted && !controllerutil.ContainsFinalizer(isvc, r.modelRegistryFinalizer) {
		log.Info("Adding Finalizer for ModelRegistry")

		if ok := controllerutil.AddFinalizer(isvc, r.modelRegistryFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the InferenceService custom resource")

			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.client.Update(ctx, isvc); err != nil {
			log.Error(err, "Failed to update InferenceService custom resource to add finalizer")

			return ctrl.Result{}, err
		}
	}

	// Retrieve or create the ServingEnvironment associated to the current namespace
	servingEnvironment, err := r.getOrCreateServingEnvironment(mrApiCtx, log, mrApi, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	if okMrIsvcId {
		// Retrieve the IS from model registry using the id
		log.Info("Retrieving model registry InferenceService by id", "mrIsvcId", mrIsvcId)
		mrIs, _, err = mrApi.ModelRegistryServiceAPI.GetInferenceService(mrApiCtx, mrIsvcId).Execute()
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("unable to find InferenceService with id %s in model registry: %w", mrIsvcId, err)
		}
	} else if okRegisteredModelId {
		// No corresponding InferenceService in model registry, create new one
		mrIs, err = r.createMRInferenceService(mrApiCtx, log, mrApi, isvc, *servingEnvironment.Id, registeredModelId, modelVersionId)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if mrIs == nil {
		// This should NOT happen
		return ctrl.Result{}, fmt.Errorf("unexpected nil model registry InferenceService")
	}

	if isMarkedToBeDeleted {
		err := r.onDeletion(ctx, mrApi, log, mrIs)
		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}

		if controllerutil.ContainsFinalizer(isvc, r.modelRegistryFinalizer) {
			log.Info("Removing Finalizer for modelRegistry after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(isvc, r.modelRegistryFinalizer); !ok {
				log.Error(err, "Failed to remove modelRegistry finalizer for InferenceService")
				return ctrl.Result{Requeue: true}, nil
			}

			if err = r.client.Update(ctx, isvc); IgnoreDeletingErrors(err) != nil {
				log.Error(err, "Failed to remove modelRegistry finalizer for InferenceService")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// No need to update the ISVC, the IS id is already set
	if isvc.Labels[r.inferenceServiceIDLabel] == *mrIs.Id {
		return ctrl.Result{}, nil
	}

	// Update the ISVC label, set the newly created IS id if not present yet
	desired := isvc.DeepCopy()

	desired.Labels[r.inferenceServiceIDLabel] = *mrIs.Id

	if err := r.client.Update(ctx, desired); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InferenceServiceController) SetupWithManager(mgr ctrl.Manager) error {
	builder := ctrl.NewControllerManagedBy(mgr).
		For(&kservev1beta1.InferenceService{})

	return builder.Complete(r)
}

func (r *InferenceServiceController) initModelRegistryService(ctx context.Context, log logr.Logger, name, namespace, url string) (*openapi.APIClient, error) {
	var err error

	log1 := log.WithValues("mr-namespace", namespace, "mr-name", name)

	if url == "" {
		log1.Info("Retrieving api http port from deployed model registry service")

		url, err = r.getMRUrlFromService(ctx, name, namespace)
		if err != nil {
			log1.Error(err, "Unable to fetch the Model Registry Service")
			return nil, err
		}
	}

	cfg := &openapi.Configuration{
		HTTPClient: r.httpClient,
		Servers: openapi.ServerConfigurations{
			{
				URL: url,
			},
		},
	}

	client := openapi.NewAPIClient(cfg)

	return client, nil
}

func (r *InferenceServiceController) getMRUrlFromService(ctx context.Context, name, namespace string) (string, error) {
	svc := &corev1.Service{}

	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, svc)
	if err != nil {
		return "", err
	}

	var restApiPort *int32

	for _, port := range svc.Spec.Ports {
		if port.Name == "http-api" {
			restApiPort = &port.Port
			break
		}
	}

	if restApiPort == nil {
		return "", fmt.Errorf("unable to find the http port in the Model Registry Service")
	}

	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", name, namespace, *restApiPort), nil
}

func (r *InferenceServiceController) createMRInferenceService(
	ctx context.Context,
	log logr.Logger,
	mr *openapi.APIClient,
	isvc *kservev1beta1.InferenceService,
	servingEnvironmentId string,
	registeredModelId string,
	modelVersionId string,
) (*openapi.InferenceService, error) {
	modelVersionIdPtr := &modelVersionId
	if modelVersionId == "" {
		modelVersionIdPtr = nil
	}

	isName := fmt.Sprintf("%s/%s", isvc.Name, isvc.UID)

	is, _, err := mr.ModelRegistryServiceAPI.FindInferenceService(ctx).
		Name(isName).ParentResourceId(servingEnvironmentId).Execute()
	if err != nil {
		log.Info("Creating new model registry InferenceService", "name", isName, "registeredModelId", registeredModelId, "modelVersionId", modelVersionId)

		is, _, err = mr.ModelRegistryServiceAPI.CreateInferenceService(ctx).InferenceServiceCreate(openapi.InferenceServiceCreate{
			DesiredState:         openapi.INFERENCESERVICESTATE_DEPLOYED.Ptr(),
			ModelVersionId:       modelVersionIdPtr,
			Name:                 &isName,
			RegisteredModelId:    registeredModelId,
			Runtime:              isvc.Spec.Predictor.Model.Runtime,
			ServingEnvironmentId: servingEnvironmentId,
		}).Execute()
	}

	return is, err
}

func (r *InferenceServiceController) getOrCreateServingEnvironment(ctx context.Context, log logr.Logger, mr *openapi.APIClient, namespace string) (*openapi.ServingEnvironment, error) {
	servingEnvironment, _, err := mr.ModelRegistryServiceAPI.FindServingEnvironment(ctx).Name(namespace).Execute()
	if err != nil {
		log.Info("ServingEnvironment not found, creating it..")

		servingEnvironment, _, err = mr.ModelRegistryServiceAPI.CreateServingEnvironment(ctx).ServingEnvironmentCreate(openapi.ServingEnvironmentCreate{
			Name: &namespace,
		}).Execute()
		if err != nil {
			return nil, fmt.Errorf("unable to create ServingEnvironment: %w", err)
		}
	}

	return servingEnvironment, nil
}

// onDeletion mark model registry inference service to UNDEPLOYED desired state
func (r *InferenceServiceController) onDeletion(ctx context.Context, mr *openapi.APIClient, log logr.Logger, is *openapi.InferenceService) (err error) {
	log.Info("Running onDeletion logic")
	if is.DesiredState != nil && *is.DesiredState != openapi.INFERENCESERVICESTATE_UNDEPLOYED {
		log.Info("InferenceService going to be deleted from cluster, setting desired state to UNDEPLOYED in model registry")

		_, _, err = mr.ModelRegistryServiceAPI.UpdateInferenceService(ctx, *is.Id).InferenceServiceUpdate(openapi.InferenceServiceUpdate{
			DesiredState: openapi.INFERENCESERVICESTATE_UNDEPLOYED.Ptr(),
		}).Execute()
	}
	return err
}

func IgnoreDeletingErrors(err error) error {
	if err == nil {
		return nil
	}
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
