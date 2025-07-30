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
      const first = parseInt(a.lastUpdateTimeSinceEpoch);
      const second = parseInt(b.lastUpdateTimeSinceEpoch);
      return new Date(second).getTime() - new Date(first).getTime();
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
