package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CanNamespaceAccessRegistry(
	ctx context.Context,
	client kubernetes.Interface,
	logger *slog.Logger,
	jobNamespace, registryName, registryNamespace string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	saSubject := "system:serviceaccount:" + jobNamespace + ":default"
	sar := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			User: saSubject,
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:      "get",
				Resource:  "services",
				Namespace: registryNamespace,
				Name:      registryName,
			},
		},
	}

	response, err := client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf("SAR failed: %w", err)
	}
	if !response.Status.Allowed {
		logger.Warn("access denied for namespace registry access",
			"jobNamespace", jobNamespace,
			"registry", registryName,
			"registryNamespace", registryNamespace,
		)
		return false, nil
	}
	return true, nil
}
