import { ModelRegistryCustomProperties, ModelState, RegisteredModel } from '~/app/types';

type MockRegisteredModelType = {
  id?: string;
  name?: string;
  owner?: string;
  state?: ModelState;
  description?: string;
  customProperties?: ModelRegistryCustomProperties;
};

export const mockRegisteredModel = ({
  name = 'test',
  owner = 'Author 1',
  state = ModelState.LIVE,
  description = '',
  customProperties = {},
  id = '1',
}: MockRegisteredModelType): RegisteredModel => ({
  createTimeSinceEpoch: '1710404288975',
  description,
  externalID: '1234132asdfasdf',
  id,
  lastUpdateTimeSinceEpoch: '1710404288975',
  name,
  state,
  owner,
  customProperties,
});
