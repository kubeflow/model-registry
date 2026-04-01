import { ModelTransferJobStatus } from '~/app/types';
import { getStatusLabel } from '~/app/pages/modelRegistry/screens/ModelTransferJobs/ModelTransferJobTableRow';

describe('getStatusLabel', () => {
  it('returns green "Complete" for COMPLETED status', () => {
    const result = getStatusLabel(ModelTransferJobStatus.COMPLETED);
    expect(result.label).toBe('Complete');
    expect(result.color).toBe('green');
  });

  it('returns blue "Running" for RUNNING status', () => {
    const result = getStatusLabel(ModelTransferJobStatus.RUNNING);
    expect(result.label).toBe('Running');
    expect(result.color).toBe('blue');
  });

  it('returns purple "Pending" for PENDING status', () => {
    const result = getStatusLabel(ModelTransferJobStatus.PENDING);
    expect(result.label).toBe('Pending');
    expect(result.color).toBe('purple');
  });

  it('returns red "Failed" for FAILED status', () => {
    const result = getStatusLabel(ModelTransferJobStatus.FAILED);
    expect(result.label).toBe('Failed');
    expect(result.color).toBe('red');
  });

  it('returns grey "Canceled" for CANCELLED status', () => {
    const result = getStatusLabel(ModelTransferJobStatus.CANCELLED);
    expect(result.label).toBe('Canceled');
    expect(result.color).toBe('grey');
  });
});
