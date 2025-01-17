import { APIOptions } from '~/shared/api/types';
import { handleRestFailures } from '~/shared/api/errorUtils';
import { isModelRegistryResponse, restGET } from '~/shared/api/apiUtils';
import { ModelRegistry } from '~/app/types';
import { BFF_API_VERSION } from '~/app/const';
import { URL_PREFIX } from '~/shared/utilities/const';
import { Namespace, UserSettings } from '~/shared/types';

export const getListModelRegistries =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<ModelRegistry[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/model_registry`, queryParams, opts),
    ).then((response) => {
      if (isModelRegistryResponse<ModelRegistry[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getUser =
  (hostPath: string) =>
  (opts: APIOptions): Promise<UserSettings> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/user`, {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<UserSettings>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getNamespaces =
  (hostPath: string) =>
  (opts: APIOptions): Promise<Namespace[]> =>
    handleRestFailures(
      restGET(hostPath, `${URL_PREFIX}/api/${BFF_API_VERSION}/namespaces`, {}, opts),
    ).then((response) => {
      if (isModelRegistryResponse<Namespace[]>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });
