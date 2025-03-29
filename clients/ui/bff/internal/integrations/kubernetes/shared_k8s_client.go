package kubernetes

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log/slog"
	"time"
)

type SharedClientLogic struct {
	Client kubernetes.Interface
	Logger *slog.Logger
	Token  BearerToken
}

func (kc *SharedClientLogic) GetServiceNames(sessionCtx context.Context, namespace string) ([]string, error) {
	services, err := kc.GetServiceDetails(sessionCtx, namespace)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(services))
	for _, svc := range services {
		names = append(names, svc.Name)
	}

	return names, nil
}

func (kc *SharedClientLogic) GetServiceDetails(sessionCtx context.Context, namespace string) ([]ServiceDetails, error) {

	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessionLogger := sessionCtx.Value(constants.TraceLoggerKey).(*slog.Logger)

	labelSelector := fmt.Sprintf("component=%s", ComponentLabelValue)

	serviceList, err := kc.Client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ServiceDetails

	for _, service := range serviceList.Items {
		var httpPort int32
		hasHTTPPort := false
		for _, port := range service.Spec.Ports {
			if port.Name == "http-api" {
				httpPort = port.Port
				hasHTTPPort = true
				break
			}
		}
		if !hasHTTPPort {
			sessionLogger.Error("service missing HTTP port", "serviceName", service.Name)
			continue
		}

		if service.Spec.ClusterIP == "" {
			sessionLogger.Error("service missing valid ClusterIP", "serviceName", service.Name)
			continue
		}

		displayName := ""
		description := ""

		if service.Annotations != nil {
			displayName = service.Annotations["displayName"]
			description = service.Annotations["description"]
		}

		if displayName == "" {
			sessionLogger.Warn("service missing displayName annotation", "serviceName", service.Name)
		}

		if description == "" {
			sessionLogger.Warn("service missing description annotation", "serviceName", service.Name)
		}

		serviceDetails := ServiceDetails{
			Name:        service.Name,
			DisplayName: displayName,
			Description: description,
			ClusterIP:   service.Spec.ClusterIP,
			HTTPPort:    httpPort,
		}

		services = append(services, serviceDetails)

	}

	return services, nil
}

func (kc *SharedClientLogic) GetServiceDetailsByName(sessionCtx context.Context, namespace string, serviceName string) (ServiceDetails, error) {
	services, err := kc.GetServiceDetails(sessionCtx, namespace)
	if err != nil {
		return ServiceDetails{}, fmt.Errorf("failed to get service details: %w", err)
	}

	for _, service := range services {
		if service.Name == serviceName {
			return service, nil
		}
	}

	return ServiceDetails{}, fmt.Errorf("service %s not found", serviceName)
}

func (kc *SharedClientLogic) BearerToken() (string, error) {
	// Token is retained for follow-up calls; do not log it.
	return kc.Token.Raw(), nil
}
