import useFetch, { FetchState } from '~/app/utils/useFetch';
import { listModelRegistryCertificateNames } from '~/app/services/modelRegistrySettingsService';
import * as React from 'react';
import { ListConfigSecretsResponse } from '~/app/k8sTypes';

const useModelRegistryCertificateNames = (
  shouldFetch: boolean,
): FetchState<ListConfigSecretsResponse> => {
  const getCertificates = React.useCallback(() => {
    if (!shouldFetch) {
      return Promise.resolve({ secrets: [], configMaps: [] });
    }
    return listModelRegistryCertificateNames();
  }, [shouldFetch]);

  return useFetch<ListConfigSecretsResponse>(
    getCertificates,
    { secrets: [], configMaps: [] },
    { refreshRate: 0 },
  );
};

export default useModelRegistryCertificateNames; 