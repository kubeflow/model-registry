import { kebabTableColumn, SortableData, ModelRegistryKind } from 'mod-arch-shared';

export const modelRegistryColumns: SortableData<ModelRegistryKind>[] = [
  {
    field: 'model regisry name',
    label: 'Model registry name',
    sortable: (a, b) => a.metadata.name.localeCompare(b.metadata.name),
    width: 30,
  },
  {
    field: 'status',
    label: 'Status',
    sortable: false,
  },
  {
    field: 'manage permissions',
    label: '',
    sortable: false,
  },
  kebabTableColumn(),
];
