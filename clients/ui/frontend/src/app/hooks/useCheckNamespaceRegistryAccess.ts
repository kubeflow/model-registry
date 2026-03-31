import * as React from 'react';
import { checkNamespaceRegistryAccess } from '~/app/api/k8s';

export type UseCheckNamespaceRegistryAccessResult = {
  hasAccess: boolean | undefined;
  isLoading: boolean;
  error: Error | undefined;
  cannotCheck: boolean;
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
  const [cannotCheck, setCannotCheck] = React.useState(false);

  React.useEffect(() => {
    if (!jobNamespace || !registryName || !registryNamespace) {
      setHasAccess(undefined);
      setError(undefined);
      setCannotCheck(false);
      return;
    }

    let cancelled = false;
    setIsLoading(true);
    setError(undefined);
    setHasAccess(undefined);
    setCannotCheck(false);

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
          setCannotCheck(result.cannotCheck);
          setHasAccess(result.cannotCheck ? undefined : result.hasAccess);
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

  return { hasAccess, isLoading, error, cannotCheck };
};
