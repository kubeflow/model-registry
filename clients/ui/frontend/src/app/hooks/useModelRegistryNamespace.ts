import { useQueryParamNamespaces } from 'mod-arch-core';
import { MODEL_REGISTRY_NAMESPACE } from '~/app/utilities/const';

/**
 * Returns the namespace where the model registry is deployed.
 * Used for SSAR validation (check-namespace-registry-access) and register-and-store flows.
 *
 * Priority:
 * 1. MODEL_REGISTRY_NAMESPACE env (set by distribution: e.g. RHOAI=rhoai-model-registries, ODH=odh-model-registries)
 * 2. namespace from URL query params (from dashboard/central namespace selector)
 *
 * Downstream distributions can set MODEL_REGISTRY_NAMESPACE at build/runtime or replace this
 * with a DSC (DataScienceCluster) status–based hook if available.
 */
export function useModelRegistryNamespace(): string | undefined {
  const queryParams = useQueryParamNamespaces();
  const fromQuery = typeof queryParams.namespace === 'string' ? queryParams.namespace : undefined;
  return MODEL_REGISTRY_NAMESPACE ?? fromQuery;
}
