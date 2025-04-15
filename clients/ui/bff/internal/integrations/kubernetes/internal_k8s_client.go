package kubernetes

import (
	"context"
	"fmt"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log/slog"
	"time"
)

type InternalKubernetesClient struct {
	SharedClientLogic
}

// newInternalKubernetesClient creates a Kubernetes client
// using the credentials of the running backend to create a single instance of the client
// If running inside the cluster, it uses the pod's service account.
// If running locally (e.g. for development), it uses the current user's kubeconfig context.
func newInternalKubernetesClient(logger *slog.Logger) (KubernetesClientInterface, error) {
	// Get kubeconfig
	kubeconfig, err := helper.GetKubeconfig()
	if err != nil {
		logger.Error("failed to get kubeconfig", "error", err)
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	// Create client
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		logger.Error("failed to create Kubernetes client", "error", err)
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &InternalKubernetesClient{
		SharedClientLogic: SharedClientLogic{
			Client: clientset,
			Logger: logger,
			Token:  NewBearerToken(kubeconfig.BearerToken),
		},
	}, nil
}

func (kc *InternalKubernetesClient) CanListServicesInNamespace(ctx context.Context, identity *RequestIdentity, namespace string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Perform SAR for get and list verbs
	for _, verb := range []string{"get", "list"} {
		sar := &authv1.SubjectAccessReview{
			Spec: authv1.SubjectAccessReviewSpec{
				User:   identity.UserID,
				Groups: identity.Groups,
				ResourceAttributes: &authv1.ResourceAttributes{
					Verb:      verb,
					Resource:  "services",
					Namespace: namespace,
				},
			},
		}

		response, err := kc.Client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			return false, fmt.Errorf("SAR failed: %w", err)
		}

		if !response.Status.Allowed {
			kc.Logger.Warn("access denied", "user", identity.UserID, "verb", verb, "resource", "services", "namespace", namespace)
			return false, nil
		}
	}

	return true, nil
}

func (kc *InternalKubernetesClient) CanAccessServiceInNamespace(ctx context.Context, identity *RequestIdentity, namespace, serviceName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sar := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			User:   identity.UserID,
			Groups: identity.Groups,
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:      "get",
				Resource:  "services",
				Namespace: namespace,
				Name:      serviceName,
			},
		},
	}

	// Perform SAR
	response, err := kc.Client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf("SAR failed: %w", err)
	}

	if !response.Status.Allowed {
		kc.Logger.Warn("access denied",
			"user", identity.UserID,
			"verb", "get",
			"resource", "services",
			"namespace", namespace,
			"service", serviceName,
		)
		return false, nil
	}

	return true, nil
}

func (kc *InternalKubernetesClient) GetNamespaces(ctx context.Context, identity *RequestIdentity) ([]corev1.Namespace, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	//list namespaces
	namespaceList, err := kc.Client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	//check access for each namespace
	var allowed []corev1.Namespace
	for _, ns := range namespaceList.Items {
		sar := &authv1.SubjectAccessReview{
			Spec: authv1.SubjectAccessReviewSpec{
				User:   identity.UserID,
				Groups: identity.Groups,
				ResourceAttributes: &authv1.ResourceAttributes{
					Verb:      "get",
					Resource:  "namespaces",
					Namespace: ns.Name,
				},
			},
		}

		response, err := kc.Client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			kc.Logger.Error("failed SAR for namespace", "namespace", ns.Name, "error", err)
			continue
		}

		if response.Status.Allowed {
			allowed = append(allowed, ns)
		}
	}

	return allowed, nil
}

func (kc *InternalKubernetesClient) IsClusterAdmin(identity *RequestIdentity) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	crbList, err := kc.Client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		kc.Logger.Error("failed to list ClusterRoleBindings", "error", err)
		return false, fmt.Errorf("failed to list ClusterRoleBindings: %w", err)
	}

	for _, crb := range crbList.Items {
		if crb.RoleRef.Kind != "ClusterRole" || crb.RoleRef.Name != "cluster-admin" {
			continue
		}
		for _, subject := range crb.Subjects {
			if subject.Kind == "User" && subject.Name == identity.UserID {
				kc.Logger.Info("user is cluster-admin", "user", identity.UserID, "crb", crb.Name)
				return true, nil
			}
		}
	}

	kc.Logger.Info("user is not cluster-admin", "user", identity.UserID)
	return false, nil
}

func (kc *InternalKubernetesClient) GetUser(identity *RequestIdentity) (string, error) {
	// On internal client, we can use the identity from request directly
	return identity.UserID, nil
}
