import { ModelTransferJobStatus } from '~/app/types';
import { getStatusLabel } from '~/app/pages/modelRegistry/screens/ModelTransferJobs/ModelTransferJobTableRow';

describe('getStatusLabel', () => {
  it.each([
    [ModelTransferJobStatus.COMPLETED, 'Complete', 'green'],
    [ModelTransferJobStatus.RUNNING, 'Running', 'blue'],
    [ModelTransferJobStatus.PENDING, 'Pending', 'purple'],
    [ModelTransferJobStatus.FAILED, 'Failed', 'red'],
    [ModelTransferJobStatus.CANCELLED, 'Canceled', 'grey'],
  ])(
    'returns correct label, color, and icon for %s status',
    (status, expectedLabel, expectedColor) => {
      const result = getStatusLabel(status);
      expect(result.label).toBe(expectedLabel);
      expect(result.color).toBe(expectedColor);
      expect(result.icon).toBeDefined();
    },
  );
});
