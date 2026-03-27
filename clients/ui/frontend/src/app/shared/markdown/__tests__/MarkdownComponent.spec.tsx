import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import ReactMarkdown from 'react-markdown';
import rehypeUnwrapImages from 'rehype-unwrap-images';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import rehypeRaw from 'rehype-raw';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';

jest.mock('react-markdown', () => ({
  __esModule: true,
  default: jest.fn(({ children }) => {
    // Basic mock implementation for headings
    if (typeof children === 'string' && children.startsWith('# ')) {
      return <h1>{children.substring(2)}</h1>;
    }
    return <>{children}</>;
  }),
}));

jest.mock('rehype-unwrap-images', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('remark-gfm', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('rehype-sanitize', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('rehype-raw', () => ({
  __esModule: true,
  default: jest.fn(),
}));

describe('MarkdownComponent', () => {
  beforeEach(() => {
    (ReactMarkdown as jest.Mock).mockClear();
    (rehypeUnwrapImages as jest.Mock).mockClear();
    (remarkGfm as jest.Mock).mockClear();
    (rehypeSanitize as jest.Mock).mockClear();
    (rehypeRaw as jest.Mock).mockClear();
  });

  it('renders markdown content', () => {
    render(<MarkdownComponent data="# Hello" />);
    expect(screen.getByRole('heading', { name: /hello/i })).toBeInTheDocument();
  });

  it('passes markdown content and plugins to ReactMarkdown', () => {
    const markdown = `# Heading\n\nThis is a [link](https://example.com).`;
    render(<MarkdownComponent data={markdown} dataTestId="markdown" />);

    expect(ReactMarkdown).toHaveBeenCalledTimes(1);

    const receivedProps = (ReactMarkdown as jest.Mock).mock.calls[0][0];

    expect(receivedProps.children).toBe(markdown);
    expect(receivedProps.remarkPlugins).toEqual([remarkGfm]);
    expect(receivedProps.rehypePlugins).toEqual(
      expect.arrayContaining([rehypeRaw, rehypeSanitize, rehypeUnwrapImages]),
    );
    expect(receivedProps.components).toBeDefined();
    expect(typeof receivedProps.components).toBe('object');
  });

  it('provides pre and code component mappers', () => {
    render(<MarkdownComponent data="test" />);

    const receivedProps = (ReactMarkdown as jest.Mock).mock.calls[0][0];
    const { pre, code } = receivedProps.components;

    expect(pre).toBeDefined();
    expect(code).toBeDefined();
  });

  it('code mapper renders inline code element', () => {
    render(<MarkdownComponent data="test" />);

    const receivedProps = (ReactMarkdown as jest.Mock).mock.calls[0][0];
    const CodeMapper = receivedProps.components.code;

    const { container } = render(<CodeMapper className="language-js">const x = 1;</CodeMapper>);

    const codeEl = container.querySelector('code');
    expect(codeEl).toBeInTheDocument();
    expect(codeEl).toHaveClass('language-js');
    expect(codeEl).toHaveTextContent('const x = 1;');
  });

  it('pre mapper extracts text from code children', () => {
    render(<MarkdownComponent data="test" />);

    const receivedProps = (ReactMarkdown as jest.Mock).mock.calls[0][0];
    const PreMapper = receivedProps.components.pre;

    const { container } = render(
      <PreMapper>
        <code>pip install torch</code>
      </PreMapper>,
    );

    expect(container.querySelector('.pf-v6-c-code-block')).toBeInTheDocument();
    expect(screen.getByText('pip install torch')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Copy to clipboard' })).toBeInTheDocument();
  });
});
