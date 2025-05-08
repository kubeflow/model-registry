import { kebabTableColumn, SortableData, isPlatformDefault } from 'mod-arch-shared';
import { ModelRegistry } from '~/app/types';

export const modelRegistryColumns: SortableData<ModelRegistry>[] = [
  {
    field: 'model regisry name',
    label: 'Model registry name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 30,
  },
  {
    field: 'status',
    label: 'Status',
    sortable: false,
  },
  ...(isPlatformDefault()
    ? [
        {
          field: 'manage permissions',
          label: '',
          sortable: false,
        },
        kebabTableColumn(),
      ]
    : []),
];
