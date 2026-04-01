import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import '@testing-library/jest-dom';
import JobNamespaceInput from '~/app/pages/modelRegistry/screens/ModelTransferJobs/JobNamespaceInput';

describe('JobNamespaceInput', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('renders with the correct label and helper text', () => {
    render(<JobNamespaceInput value="" onChange={jest.fn()} />);
    expect(screen.getByText('Transfer job namespace')).toBeInTheDocument();
    expect(
      screen.getByText(/you do not have permission to list transfer jobs/i),
    ).toBeInTheDocument();
  });

  it('does not call onChange immediately on typing', () => {
    const onChange = jest.fn();
    render(<JobNamespaceInput value="" onChange={onChange} />);

    fireEvent.change(screen.getByTestId('job-namespace-input'), {
      target: { value: 'my-namespace' },
    });

    expect(onChange).not.toHaveBeenCalled();
  });

  it('calls onChange after 2-second debounce', () => {
    const onChange = jest.fn();
    render(<JobNamespaceInput value="" onChange={onChange} />);

    fireEvent.change(screen.getByTestId('job-namespace-input'), {
      target: { value: 'my-namespace' },
    });

    act(() => {
      jest.advanceTimersByTime(2000);
    });

    expect(onChange).toHaveBeenCalledWith('my-namespace');
    expect(onChange).toHaveBeenCalledTimes(1);
  });

  it('calls onChange on blur before debounce fires', () => {
    const onChange = jest.fn();
    render(<JobNamespaceInput value="" onChange={onChange} />);
    const input = screen.getByTestId('job-namespace-input');

    fireEvent.change(input, { target: { value: 'blur-namespace' } });
    fireEvent.blur(input);

    expect(onChange).toHaveBeenCalledWith('blur-namespace');
    expect(onChange).toHaveBeenCalledTimes(1);
  });

  it('resets debounce timer on subsequent keystrokes', () => {
    const onChange = jest.fn();
    render(<JobNamespaceInput value="" onChange={onChange} />);
    const input = screen.getByTestId('job-namespace-input');

    fireEvent.change(input, { target: { value: 'first' } });

    act(() => {
      jest.advanceTimersByTime(1500);
    });

    fireEvent.change(input, { target: { value: 'second' } });

    act(() => {
      jest.advanceTimersByTime(1500);
    });

    // First debounce should have been cancelled
    expect(onChange).not.toHaveBeenCalled();

    act(() => {
      jest.advanceTimersByTime(500);
    });

    expect(onChange).toHaveBeenCalledWith('second');
    expect(onChange).toHaveBeenCalledTimes(1);
  });
});
