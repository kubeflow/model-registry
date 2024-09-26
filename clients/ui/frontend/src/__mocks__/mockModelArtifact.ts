import { ModelArtifact } from '~/app/types';

export const mockModelArtifact = (partial?: Partial<ModelArtifact>): ModelArtifact => ({
  createTimeSinceEpoch: '1712234877179',
  id: '1',
  lastUpdateTimeSinceEpoch: '1712234877179',
  name: 'fraud detection model version 1',
  description: 'Description of model version',
  artifactType: 'model-artifact',
  customProperties: {},
  storageKey: 'test storage key',
  storagePath: 'test path',
  uri: 's3://test-bucket/demo-models/test-path?endpoint=test-endpoint&defaultRegion=test-region',
  modelFormatName: 'test model format',
  modelFormatVersion: 'test version 1',
  ...partial,
});
