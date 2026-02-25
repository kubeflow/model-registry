import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import {
  ModelTransferJob,
  ModelTransferJobStatus,
  ModelTransferJobUploadIntent,
} from '~/app/types';
import RetryJobModal from '~/app/pages/modelRegistry/screens/ModelTransferJobs/RetryJobModal';

describe('RetryJobModal', () => {
  const mockOnClose = jest.fn();
  const mockOnRetry = jest.fn();

  const mockJob: ModelTransferJob = {
    id: 'test-job-id',
    name: 'test-job-name',
    status: ModelTransferJobStatus.FAILED,
    uploadIntent: ModelTransferJobUploadIntent.CREATE_MODEL,
    namespace: 'test-namespace',
    registeredModelName: 'test-model',
    modelVersionName: 'v1.0.0',
    createTimeSinceEpoch: Date.now().toString(),
  };

  beforeEach(() => {
    mockOnClose.mockClear();
    mockOnRetry.mockClear();
    mockOnRetry.mockResolvedValue(undefined);
  });

  it('should render the modal with correct title', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    expect(screen.getByText('Retry model transfer job?')).toBeInTheDocument();
  });

  it('should display the model version name in the description', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    expect(screen.getByText('v1.0.0')).toBeInTheDocument();
  });

  it('should auto-generate a retry job name with -2 suffix', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const input = screen.getByTestId('retry-job-name-input');
    expect(input).toHaveValue('test-job-name-2');
  });

  it('should increment existing numeric suffix', () => {
    const jobWithSuffix: ModelTransferJob = {
      ...mockJob,
      name: 'my-job-3',
    };

    render(<RetryJobModal job={jobWithSuffix} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const input = screen.getByTestId('retry-job-name-input');
    expect(input).toHaveValue('my-job-4');
  });

  it('should have delete checkbox checked by default', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const checkbox = screen.getByTestId('delete-old-job-checkbox');
    expect(checkbox).toBeChecked();
  });

  it('should display the old job name in the delete checkbox label', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    expect(screen.getByText('test-job-name')).toBeInTheDocument();
  });

  it('should call onClose when Cancel button is clicked', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    fireEvent.click(screen.getByText('Cancel'));
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('should call onRetry with correct parameters when Retry button is clicked', async () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    fireEvent.click(screen.getByTestId('retry-job-submit-button'));

    await waitFor(() => {
      expect(mockOnRetry).toHaveBeenCalledWith('test-job-name-2', true);
    });
  });

  it('should call onRetry with deleteOldJob=false when checkbox is unchecked', async () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const checkbox = screen.getByTestId('delete-old-job-checkbox');
    fireEvent.click(checkbox);

    fireEvent.click(screen.getByTestId('retry-job-submit-button'));

    await waitFor(() => {
      expect(mockOnRetry).toHaveBeenCalledWith('test-job-name-2', false);
    });
  });

  it('should allow editing the job name', async () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const input = screen.getByTestId('retry-job-name-input');
    fireEvent.change(input, { target: { value: 'custom-retry-name' } });

    fireEvent.click(screen.getByTestId('retry-job-submit-button'));

    await waitFor(() => {
      expect(mockOnRetry).toHaveBeenCalledWith('custom-retry-name', true);
    });
  });

  it('should disable Retry button when job name is empty', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const input = screen.getByTestId('retry-job-name-input');
    fireEvent.change(input, { target: { value: '' } });

    const retryButton = screen.getByTestId('retry-job-submit-button');
    expect(retryButton).toBeDisabled();
  });

  it('should disable Retry button when job name has invalid characters', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const input = screen.getByTestId('retry-job-name-input');
    fireEvent.change(input, { target: { value: 'Invalid_Name!' } });

    const retryButton = screen.getByTestId('retry-job-submit-button');
    expect(retryButton).toBeDisabled();
  });

  it('should show error message for invalid job name', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    const input = screen.getByTestId('retry-job-name-input');
    fireEvent.change(input, { target: { value: 'Invalid_Name!' } });

    expect(
      screen.getByText(/Must start and end with a lowercase letter or number/),
    ).toBeInTheDocument();
  });

  it('should display error alert when retry fails', async () => {
    const errorMessage = 'Failed to create retry job';
    mockOnRetry.mockRejectedValue(new Error(errorMessage));

    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    fireEvent.click(screen.getByTestId('retry-job-submit-button'));

    await waitFor(() => {
      expect(screen.getByTestId('retry-job-error-alert')).toBeInTheDocument();
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  it('should show loading state while retrying', async () => {
    mockOnRetry.mockImplementation(
      () =>
        new Promise((resolve) => {
          setTimeout(resolve, 1000);
        }),
    );

    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    fireEvent.click(screen.getByTestId('retry-job-submit-button'));

    // Button should be disabled while loading
    expect(screen.getByTestId('retry-job-submit-button')).toBeDisabled();
  });

  it('should close modal after successful retry', async () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    fireEvent.click(screen.getByTestId('retry-job-submit-button'));

    await waitFor(() => {
      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });
  });

  it('should show Edit resource name link initially', () => {
    render(<RetryJobModal job={mockJob} onClose={mockOnClose} onRetry={mockOnRetry} />);

    expect(screen.getByTestId('retry-job-edit-resource-link')).toBeInTheDocument();
  });
});
