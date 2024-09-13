/* eslint-disable camelcase */
import { ModelVersionList } from '~/app/types';

export const mockModelVersionList = ({
  items = [],
}: Partial<ModelVersionList>): ModelVersionList => ({
  items,
  nextPageToken: '',
  pageSize: 0,
  size: items.length,
});
