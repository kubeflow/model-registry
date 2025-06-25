import React from 'react';
import { ModelRegistryKind, APIOptions } from 'mod-arch-shared';
import { getModelRegistrySettings } from '~/app/api/k8s';

const useModelRegistryCR = (
  namespace: string | undefined,
  name: string,
): [ModelRegistryKind | null, boolean, Error | undefined] => {
  const [modelRegistry, setModelRegistry] = React.useState<ModelRegistryKind | null>(null);
  const [loaded, setLoaded] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>(undefined);

  React.useEffect(() => {
    if (!name) {
      setModelRegistry(null);
      setLoaded(true);
      setError(undefined);
      return;
    }

    const fetchModelRegistry = async () => {
      try {
        setLoaded(false);
        setError(undefined);
        const opts: APIOptions = {};
        const hostPath = window.location.origin;
        const result = await getModelRegistrySettings(hostPath, {})(opts, name);
        setModelRegistry(result);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to fetch model registry'));
        setModelRegistry(null);
      } finally {
        setLoaded(true);
      }
    };

    fetchModelRegistry();
  }, [name, namespace]);

  return [modelRegistry, loaded, error];
};

export { useModelRegistryCR };
