/* eslint-disable camelcase */
import { ModelArtifactList } from '~/app/types';

export const mockModelArtifactList = ({
  items = [],
}: Partial<ModelArtifactList>): ModelArtifactList => ({
  items,
  nextPageToken: '',
  pageSize: 0,
  size: 1,
});
