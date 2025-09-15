import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import DetailsComponent from '~/app/shared/markdown/components/DetailsComponent';

describe('DetailsComponent', () => {
  it('renders in initial collapsed state', () => {
    render(<DetailsComponent summary="Test Summary">Test Content</DetailsComponent>);
    expect(screen.getByText('Test Summary')).toBeInTheDocument();
    expect(screen.queryByText('Test Content')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Test Summary/i })).toBeInTheDocument();
  });

  it('expands and collapses on button click', () => {
    render(<DetailsComponent summary="Click me">Hidden content</DetailsComponent>);
    expect(screen.queryByText('Hidden content')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /Click me/i }));
    expect(screen.getByText('Hidden content')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /Click me/i }));
    expect(screen.queryByText('Hidden content')).not.toBeInTheDocument();
  });

  it('processes summary text correctly', () => {
    const summaryWithNewlines = 'Line 1\\nLine 2\\nLine 3';
    render(<DetailsComponent summary={summaryWithNewlines}>Content</DetailsComponent>);
    expect(screen.getByText('Line 1 Line 2 Line 3')).toBeInTheDocument();
  });

  it('renders code blocks correctly', () => {
    const codeContent = 'const x = 42;';
    render(
      <DetailsComponent summary="Code Example">
        <code>{codeContent}</code>
      </DetailsComponent>,
    );
    fireEvent.click(screen.getByRole('button', { name: /Code Example/i }));
    const codeElement = screen.getByText(codeContent);
    expect(codeElement).toBeInTheDocument();
    expect(codeElement.closest('.pf-v6-c-code-block')).toBeInTheDocument();
  });

  it('handles multiple code blocks and regular content', () => {
    render(
      <DetailsComponent summary="Mixed Content">
        <p>Regular text</p>
        <code>First code block</code>
        <p>More text</p>
        <code>Second code block</code>
      </DetailsComponent>,
    );

    fireEvent.click(screen.getByRole('button', { name: /Mixed Content/i }));
    expect(screen.getByText('Regular text')).toBeInTheDocument();
    expect(screen.getByText('First code block')).toBeInTheDocument();
    expect(screen.getByText('More text')).toBeInTheDocument();
    expect(screen.getByText('Second code block')).toBeInTheDocument();
    const codeBlocks = document.querySelectorAll('.pf-v6-c-code-block');
    expect(codeBlocks).toHaveLength(2);
  });
});
