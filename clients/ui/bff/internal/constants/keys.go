package constants

type contextKey string

const (
	ModelRegistryHttpClientKey  contextKey = "ModelRegistryHttpClientKey"
	NamespaceHeaderParameterKey contextKey = "namespace"

	//Kubeflow authorization operates using custom authentication headers:
	// Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
	// but it's supported on Model Registry BFF
	KubeflowUserIdKey          contextKey = "kubeflowUserId" // kubeflow-userid :contains the user's email address
	KubeflowUserIDHeader                  = "kubeflow-userid"
	KubeflowUserGroupsKey      contextKey = "kubeflowUserGroups" // kubeflow-groups : Holds a comma-separated list of user groups
	KubeflowUserGroupsIdHeader            = "kubeflow-groups"

	TraceIdKey     contextKey = "TraceIdKey"
	TraceLoggerKey contextKey = "TraceLoggerKey"
)
