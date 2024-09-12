import { APIOptions } from '~/app/api/types';
import { handleRestFailures } from '~/app/api/errorUtils';
import { isModelRegistryResponse, restGET } from '~/app/api/apiUtils';
import { ModelRegistry } from '~/app/types';
import { BFF_API_VERSION } from '~/app/const';

export const getListModelRegistries =
  (hostPath: string) =>
  (opts: APIOptions): Promise<ModelRegistry[]> =>
    handleRestFailures(restGET(hostPath, `/api/${BFF_API_VERSION}/model_registry`, {}, opts)).then(
      (response) => {
        if (isModelRegistryResponse(response)) {
          return response.model_registry;
        }
        throw new Error('Invalid response format');
      },
    );
