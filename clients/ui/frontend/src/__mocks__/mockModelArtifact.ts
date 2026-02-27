import { ModelArtifact } from '~/app/types';
import { ModelSourceKind } from '~/concepts/modelRegistry/types';

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

/**
 * Creates a mock model artifact that was registered via a transfer job (Register + Store flow).
 * The modelSource* properties reference the transfer job that performed the registration.
 */
export const mockModelArtifactWithTransferJob = (
  partial?: Partial<ModelArtifact>,
): ModelArtifact => ({
  ...mockModelArtifact(),
  uri: 'oci://quay.io/my-org/my-model:v1.0.0',
  modelSourceKind: ModelSourceKind.TRANSFER_JOB,
  modelSourceGroup: 'my-project-1',
  modelSourceName: 'model-transfer-job-1',
  ...partial,
});
