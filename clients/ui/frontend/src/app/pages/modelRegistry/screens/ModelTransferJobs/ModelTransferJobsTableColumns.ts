import { SortableData } from 'mod-arch-shared';
import { ModelTransferJob } from '~/app/types';

export const modelTransferJobsColumns: SortableData<ModelTransferJob>[] = [
  {
    field: 'name',
    label: 'Job name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 15,
  },
  {
    field: 'modelName',
    label: 'Model name',
    sortable: (a, b) => (a.registeredModelName || '').localeCompare(b.registeredModelName || ''),
    width: 15,
  },
  {
    field: 'modelVersionName',
    label: 'Model version name',
    sortable: (a, b) => (a.modelVersionName || '').localeCompare(b.modelVersionName || ''),
    width: 15,
  },
  {
    field: 'namespace',
    label: 'Namespace',
    sortable: (a, b) => (a.namespace || '').localeCompare(b.namespace || ''),
    width: 10,
  },
  {
    field: 'created',
    label: 'Created',
    sortable: (a: ModelTransferJob, b: ModelTransferJob): number => {
      const timeA = Number(a.createTimeSinceEpoch || 0);
      const timeB = Number(b.createTimeSinceEpoch || 0);
      return timeB - timeA;
    },
    width: 10,
  },
  {
    field: 'author',
    label: 'Author',
    sortable: (a, b) => (a.author || '').localeCompare(b.author || ''),
    width: 10,
    info: {
      popover: 'The author is the user who created the transfer job.',
      popoverProps: {
        position: 'left',
      },
    },
  },
  {
    field: 'status',
    label: 'Transfer job status',
    sortable: (a, b) => a.status.localeCompare(b.status),
    width: 15,
    info: {
      popover: 'The current status of the model transfer job.',
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
