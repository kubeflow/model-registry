import { ModelVersion, ModelState } from '~/app/types';
import { createModelRegistryLabelsObject } from './utils';

type MockModelVersionType = {
  author?: string;
  id?: string;
  registeredModelId?: string;
  name?: string;
  labels?: string[];
  state?: ModelState;
  description?: string;
  createTimeSinceEpoch?: string;
  lastUpdateTimeSinceEpoch?: string;
};

export const mockModelVersion = ({
  author = 'Test author',
  registeredModelId = '1',
  name = 'new model version',
  labels = [],
  id = '1',
  state = ModelState.LIVE,
  description = 'Description of model version',
  createTimeSinceEpoch = '1712234877179',
  lastUpdateTimeSinceEpoch = '1712234877179',
}: MockModelVersionType): ModelVersion => ({
  author,
  createTimeSinceEpoch,
  customProperties: createModelRegistryLabelsObject(labels),
  id,
  lastUpdateTimeSinceEpoch,
  name,
  state,
  registeredModelId,
  description,
});
