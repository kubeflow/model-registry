import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import ContentModal, { ButtonAction } from '~/app/shared/components/modals/ContentModal';

describe('ContentModal', () => {
  const mockOnClose = jest.fn();

  beforeEach(() => {
    jest.useFakeTimers();
    mockOnClose.mockClear();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('should render a button without a variant as primary (default)', () => {
    const buttonActions: ButtonAction[] = [
      {
        label: 'No Variant Button',
        onClick: jest.fn(),
        // variant is intentionally not specified
        dataTestId: 'no-variant-button',
      },
    ];

    render(
      <ContentModal
        onClose={mockOnClose}
        title="Test Modal"
        contents={<div>Test content</div>}
        buttonActions={buttonActions}
      />,
    );

    const button = screen.getByTestId('no-variant-button');
    expect(button).toBeInTheDocument();
    expect(button).toHaveTextContent('No Variant Button');
    // PatternFly Button defaults to 'primary' variant when none is specified
    expect(button).toHaveClass('pf-m-primary');
  });

  it('should render buttons with explicit variants correctly', () => {
    const buttonActions: ButtonAction[] = [
      {
        label: 'Primary Button',
        onClick: jest.fn(),
        variant: 'primary',
        dataTestId: 'primary-button',
      },
      {
        label: 'Secondary Button',
        onClick: jest.fn(),
        variant: 'secondary',
        dataTestId: 'secondary-button',
      },
      {
        label: 'Link Button',
        onClick: jest.fn(),
        variant: 'link',
        dataTestId: 'link-button',
      },
      {
        label: 'Danger Button',
        onClick: jest.fn(),
        variant: 'danger',
        dataTestId: 'danger-button',
      },
    ];

    render(
      <ContentModal
        onClose={mockOnClose}
        title="Test Modal"
        contents={<div>Test content</div>}
        buttonActions={buttonActions}
      />,
    );

    expect(screen.getByTestId('primary-button')).toHaveClass('pf-m-primary');
    expect(screen.getByTestId('secondary-button')).toHaveClass('pf-m-secondary');
    expect(screen.getByTestId('link-button')).toHaveClass('pf-m-link');
    expect(screen.getByTestId('danger-button')).toHaveClass('pf-m-danger');
  });

  it('should render with string title', () => {
    render(
      <ContentModal
        onClose={mockOnClose}
        title="String Title"
        contents={<div>Test content</div>}
      />,
    );

    expect(screen.getByText('String Title')).toBeInTheDocument();
  });

  it('should render with ReactNode title', () => {
    render(
      <ContentModal
        onClose={mockOnClose}
        title={<span data-testid="custom-title">Custom Title Node</span>}
        contents={<div>Test content</div>}
      />,
    );

    expect(screen.getByTestId('custom-title')).toBeInTheDocument();
    expect(screen.getByText('Custom Title Node')).toBeInTheDocument();
  });

  it('should render description when provided', () => {
    render(
      <ContentModal
        onClose={mockOnClose}
        title="Test Modal"
        description="This is a description"
        contents={<div>Test content</div>}
      />,
    );

    expect(screen.getByText('This is a description')).toBeInTheDocument();
  });

  it('should render titleIconVariant when provided', () => {
    render(
      <ContentModal
        onClose={mockOnClose}
        title="Warning Modal"
        titleIconVariant="warning"
        contents={<div>Test content</div>}
      />,
    );

    // The warning icon should be rendered in the modal header
    const header = screen.getByTestId('generic-modal-header');
    expect(header).toBeInTheDocument();
  });
});
