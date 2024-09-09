import { APIOptions } from '~/app/api/types';
import { handleRestFailures } from '~/app/api/errorUtils';
import { restGET } from '~/app/api/apiUtils';
import { ModelRegistryList } from '~/app/types';
import { BFF_API_VERSION } from '~/app/const';

export const getListModelRegistries =
  (hostPath: string) =>
  (opts: APIOptions): Promise<ModelRegistryList> =>
    handleRestFailures(restGET(hostPath, `/api/${BFF_API_VERSION}/model_registry`, {}, opts));
