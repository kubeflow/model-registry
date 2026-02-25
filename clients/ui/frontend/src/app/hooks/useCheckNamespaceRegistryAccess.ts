import * as React from 'react';
import { checkNamespaceRegistryAccess } from '~/app/api/k8s';

export type UseCheckNamespaceRegistryAccessResult = {
  hasAccess: boolean | undefined;
  isLoading: boolean;
  error: Error | undefined;
};

/**
 * Checks if the selected namespace's default ServiceAccount has access to the model registry
 * (for register-and-store job validation). Runs when jobNamespace, registryName, and
 * registryNamespace are all defined.
 */
export const useCheckNamespaceRegistryAccess = (
  registryName: string | undefined,
  registryNamespace: string | undefined,
  jobNamespace: string | undefined,
): UseCheckNamespaceRegistryAccessResult => {
  const [hasAccess, setHasAccess] = React.useState<boolean | undefined>(undefined);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>(undefined);

  React.useEffect(() => {
    if (!jobNamespace || !registryName || !registryNamespace) {
      setHasAccess(undefined);
      setError(undefined);
      return;
    }

    let cancelled = false;
    setIsLoading(true);
    setError(undefined);
    setHasAccess(undefined);

    const run = async () => {
      try {
        const result = await checkNamespaceRegistryAccess('')(
          {},
          {
            namespace: jobNamespace,
            registryName,
            registryNamespace,
          },
        );
        if (!cancelled) {
          setHasAccess(result.hasAccess);
        }
      } catch (e) {
        if (!cancelled) {
          setError(e instanceof Error ? e : new Error(String(e)));
          setHasAccess(undefined);
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    run();
    return () => {
      cancelled = true;
    };
  }, [jobNamespace, registryName, registryNamespace]);

  return { hasAccess, isLoading, error };
};
