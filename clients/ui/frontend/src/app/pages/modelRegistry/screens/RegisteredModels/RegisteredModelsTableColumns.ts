import { SortableData } from '~/shared/components/table';
import { RegisteredModel } from '~/app/types';

export const rmColumns: SortableData<RegisteredModel>[] = [
  {
    field: 'model name',
    label: 'Model name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 40,
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
      tooltip: 'The owner is the user who registered the model.',
      tooltipProps: {
        isContentLeftAligned: true,
      },
    },
  },
  {
    field: 'kebab',
    label: '',
    sortable: false,
  },
];
