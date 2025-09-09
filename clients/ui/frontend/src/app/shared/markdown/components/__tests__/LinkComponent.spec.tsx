import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import LinkComponent from '~/app/shared/markdown/components/LinkComponent';

describe('LinkComponent', () => {
  const originalOpen = window.open;

  beforeEach(() => {
    window.open = jest.fn();
  });

  afterEach(() => {
    window.open = originalOpen;
    jest.clearAllMocks();
  });

  it('renders children content correctly', () => {
    render(<LinkComponent>Test Link</LinkComponent>);
    expect(screen.getByText('Test Link')).toBeInTheDocument();
  });

  it('does not open window when href is "#" (default)', () => {
    render(<LinkComponent>Default Link</LinkComponent>);
    fireEvent.click(screen.getByRole('button'));
    expect(window.open).not.toHaveBeenCalled();
  });

  it('opens link in new window with security attributes when href is provided', () => {
    const testUrl = 'https://example.com';
    render(<LinkComponent href={testUrl}>External Link</LinkComponent>);

    fireEvent.click(screen.getByRole('button'));

    expect(window.open).toHaveBeenCalledWith(testUrl, '_blank', 'noopener,noreferrer');
  });

  it('renders with nested elements', () => {
    render(
      <LinkComponent>
        <span>Nested</span>
        <strong>Content</strong>
      </LinkComponent>,
    );

    expect(screen.getByText('Nested')).toBeInTheDocument();
    expect(screen.getByText('Content')).toBeInTheDocument();
  });
});
