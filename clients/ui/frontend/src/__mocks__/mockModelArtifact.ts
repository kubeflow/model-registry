import { ModelArtifact, ModelArtifactState } from '~/app/types';

type MockModelArtifact = {
  id?: string;
  name?: string;
  uri?: string;
  state?: ModelArtifactState;
  author?: string;
};

export const mockModelArtifact = ({
  id = '1',
  name = 'test',
  uri = 'test',
  state = ModelArtifactState.LIVE,
  author = 'Author 1',
}: MockModelArtifact): ModelArtifact => ({
  id,
  name,
  externalID: '1234132asdfasdf',
  description: '',
  createTimeSinceEpoch: '1710404288975',
  lastUpdateTimeSinceEpoch: '1710404288975',
  customProperties: {},
  uri,
  state,
  author,
  modelFormatName: 'test',
  storageKey: 'test',
  storagePath: 'test',
  modelFormatVersion: 'test',
  serviceAccountName: 'test',
  artifactType: 'test',
});
