import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ManageColumnSearchInput } from '~/app/shared/components/manageColumns/ManageColumnSearchInput';

describe('ManageColumnSearchInput', () => {
  const mockOnSearch = jest.fn();

  beforeEach(() => {
    mockOnSearch.mockClear();
  });

  it('should render with default placeholder', () => {
    render(<ManageColumnSearchInput value="" onSearch={mockOnSearch} />);

    expect(screen.getByPlaceholderText('Filter by column name')).toBeInTheDocument();
  });

  it('should render with custom placeholder', () => {
    render(
      <ManageColumnSearchInput value="" onSearch={mockOnSearch} placeholder="Search metrics..." />,
    );

    expect(screen.getByPlaceholderText('Search metrics...')).toBeInTheDocument();
  });

  it('should display the provided value', () => {
    render(<ManageColumnSearchInput value="test search" onSearch={mockOnSearch} />);

    const input = screen.getByRole('textbox') as HTMLInputElement;
    expect(input.value).toBe('test search');
  });

  it('should call onSearch when typing', () => {
    render(<ManageColumnSearchInput value="" onSearch={mockOnSearch} />);

    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: 'new value' } });

    expect(mockOnSearch).toHaveBeenCalledWith('new value');
  });

  it('should call onSearch with empty string when cleared', () => {
    render(<ManageColumnSearchInput value="some text" onSearch={mockOnSearch} />);

    // Find and click the clear button
    const clearButton = screen.getByRole('button', { name: /reset/i });
    fireEvent.click(clearButton);

    expect(mockOnSearch).toHaveBeenCalledWith('');
  });

  it('should use custom dataTestId', () => {
    render(
      <ManageColumnSearchInput value="" onSearch={mockOnSearch} dataTestId="my-custom-search" />,
    );

    expect(screen.getByTestId('my-custom-search')).toBeInTheDocument();
  });

  it('should use default dataTestId', () => {
    render(<ManageColumnSearchInput value="" onSearch={mockOnSearch} />);

    expect(screen.getByTestId('manage-column-search')).toBeInTheDocument();
  });
});
