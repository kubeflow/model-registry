import { APIOptions } from '~/types';
import { handleRestFailures } from '~/app/api/errorUtils';
import { restGET } from '~/app/api/apiUtils';
import { ModelRegistry } from '~/app/types';
import { BFF_API_VERSION } from '~/app/const';

export const getModelRegistries =
  (hostPath: string) =>
  (opts: APIOptions): Promise<ModelRegistry> =>
    handleRestFailures(restGET(hostPath, `/api/${BFF_API_VERSION}/model_registry`, {}, opts));
