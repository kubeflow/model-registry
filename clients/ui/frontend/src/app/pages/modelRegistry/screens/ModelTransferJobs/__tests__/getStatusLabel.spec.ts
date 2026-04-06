import { ModelTransferJobStatus } from '~/app/types';
import { getStatusLabel } from '~/app/pages/modelRegistry/screens/ModelTransferJobs/ModelTransferJobTableRow';

describe('getStatusLabel', () => {
  it.each([
    [ModelTransferJobStatus.COMPLETED, 'Complete', undefined, 'success'],
    [ModelTransferJobStatus.RUNNING, 'Running', 'blue', undefined],
    [ModelTransferJobStatus.PENDING, 'Pending', 'purple', undefined],
    [ModelTransferJobStatus.FAILED, 'Failed', undefined, 'danger'],
    [ModelTransferJobStatus.CANCELLED, 'Canceled', 'grey', undefined],
  ])(
    'returns correct label, color/status, and icon for %s',
    (status, expectedLabel, expectedColor, expectedStatus) => {
      const result = getStatusLabel(status);
      expect(result.label).toBe(expectedLabel);
      expect(result.color).toBe(expectedColor);
      expect(result.status).toBe(expectedStatus);
      expect(result.icon).toBeDefined();
    },
  );
});
