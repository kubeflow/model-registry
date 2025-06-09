import { SortableData } from 'mod-arch-shared';
import { ModelVersion } from '~/app/types';

export const mvColumns: SortableData<ModelVersion>[] = [
  {
    field: 'version name',
    label: 'Version name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 40,
  },
  {
    field: 'last_modified',
    label: 'Last modified',
    sortable: (a: ModelVersion, b: ModelVersion): number => {
      const first = parseInt(a.lastUpdateTimeSinceEpoch);
      const second = parseInt(b.lastUpdateTimeSinceEpoch);
      return new Date(second).getTime() - new Date(first).getTime();
    },
  },
  {
    field: 'author',
    label: 'Author',
    sortable: (a: ModelVersion, b: ModelVersion): number => {
      const first = a.author || '';
      const second = b.author || '';
      return first.localeCompare(second);
    },
    info: {
      popover: 'The author is the user who registered the model version.',
    },
  },
  {
    field: 'labels',
    label: 'Labels',
    sortable: false,
    width: 35,
  },
  {
    field: 'kebab',
    label: '',
    sortable: false,
  },
];
