import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import CodeBlockComponent from '~/app/shared/markdown/components/CodeBlockComponent';

describe('CodeBlockComponent', () => {
  it('renders code content correctly', () => {
    const codeContent = 'const example = "test";';
    render(<CodeBlockComponent>{codeContent}</CodeBlockComponent>);

    expect(screen.getByText(codeContent)).toBeInTheDocument();
    expect(screen.getByText(codeContent).closest('.pf-v6-c-code-block')).toBeInTheDocument();
    expect(screen.getByText(codeContent).closest('.pf-v6-c-code-block__code')).toBeInTheDocument();
  });
});
