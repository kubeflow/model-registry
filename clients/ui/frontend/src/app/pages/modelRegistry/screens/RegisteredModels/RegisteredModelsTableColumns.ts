import { SortableData } from 'mod-arch-shared';
import { RegisteredModel } from '~/app/types';

export const rmColumns: SortableData<RegisteredModel>[] = [
  {
    field: 'model name',
    label: 'Model name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 35,
  },
  {
    field: 'latest version',
    label: 'Latest version',
    sortable: false,
    width: 15,
  },
  {
    field: 'labels',
    label: 'Labels',
    sortable: false,
    width: 35,
  },
  {
    field: 'last_modified',
    label: 'Last modified',
    sortable: (a: RegisteredModel, b: RegisteredModel): number => {
      // Convert timestamps to numbers for direct comparison
      const timeA = Number(a.lastUpdateTimeSinceEpoch || 0);
      const timeB = Number(b.lastUpdateTimeSinceEpoch || 0);
      // Sort in descending order (newest first)
      return timeB - timeA;
    },
  },
  {
    field: 'owner',
    label: 'Owner',
    sortable: true,
    info: {
      popover: 'The owner is the user who registered the model.',
      popoverProps: {
        position: 'left',
      },
    },
  },
  {
    field: 'kebab',
    label: '',
    sortable: false,
  },
];
