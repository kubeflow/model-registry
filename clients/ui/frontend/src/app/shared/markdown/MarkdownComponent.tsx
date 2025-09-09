import { Content, ContentVariants, List, ListItem } from '@patternfly/react-core';
import React from 'react';
import Markdown from 'react-markdown';
import rehypeUnwrapImages from 'rehype-unwrap-images';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import rehypeRaw from 'rehype-raw';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';

import './MarkdownComponent.scss';
import { shiftHeadingLevel } from './utils';
import LinkComponent from './components/LinkComponent';
import DetailsComponent from './components/DetailsComponent';
import CodeBlockComponent from './components/CodeBlockComponent';

type MarkdownComponentProps = {
  data: string;
  dataTestId?: string;
  maxHeading?: number;
};

const MarkdownComponent = ({
  data,
  dataTestId,
  maxHeading = 1,
}: MarkdownComponentProps): JSX.Element => (
  <div className="markdown-content" data-testid={dataTestId}>
    <Markdown
      components={{
        p: ({ children, ...props }) => (
          <Content component={ContentVariants.p} {...props}>
            {children}
          </Content>
        ),
        a: ({ children, href, ...props }) => (
          <LinkComponent href={href} {...props}>
            {children}
          </LinkComponent>
        ),
        details: ({ children, ...props }) => {
          const summary = React.Children.toArray(children).find(
            (child) => React.isValidElement(child) && child.type === 'summary',
          );
          const content = React.Children.toArray(children).filter(
            (child) => !(React.isValidElement(child) && child.type === 'summary'),
          );
          return (
            <DetailsComponent
              summary={
                summary && React.isValidElement(summary)
                  ? typeof summary.props.children === 'string'
                    ? summary.props.children
                    : 'Details'
                  : 'Details'
              }
              {...props}
            >
              {content}
            </DetailsComponent>
          );
        },
        summary: ({ children, ...props }) => <Content {...props}>{children}</Content>,
        code: ({ node, className, children, ...props }) => {
          const code = React.Children.toArray(children)
            .map((child) => (typeof child === 'string' ? child : ''))
            .join('')
            .replace(/\n$/, '');

          if (!node) {
            return (
              <code className={className} {...props}>
                {children}
              </code>
            );
          }

          const isPre = 'tagName' in node && node.tagName === 'pre';
          if (isPre) {
            return <CodeBlockComponent {...props}>{code}</CodeBlockComponent>;
          }

          return (
            <code className={className} {...props}>
              {children}
            </code>
          );
        },
        h1: ({ children, ...props }) => (
          <Content component={shiftHeadingLevel(1, maxHeading)} {...props}>
            {children}
          </Content>
        ),
        h2: ({ children, ...props }) => (
          <Content component={shiftHeadingLevel(2, maxHeading)} {...props}>
            {children}
          </Content>
        ),
        h3: ({ children, ...props }) => (
          <Content component={shiftHeadingLevel(3, maxHeading)} {...props}>
            {children}
          </Content>
        ),
        h4: ({ children, ...props }) => (
          <Content component={shiftHeadingLevel(4, maxHeading)} {...props}>
            {children}
          </Content>
        ),
        h5: ({ children, ...props }) => (
          <Content component={shiftHeadingLevel(5, maxHeading)} {...props}>
            {children}
          </Content>
        ),
        h6: ({ children, ...props }) => (
          <Content component={shiftHeadingLevel(6, maxHeading)} {...props}>
            {children}
          </Content>
        ),
        blockquote: ({ children, ...props }) => (
          <Content component={ContentVariants.blockquote} {...props}>
            {children}
          </Content>
        ),
        ul: ({ children, ...props }) => (
          <List component="ul" {...props}>
            {children}
          </List>
        ),
        ol: ({ children, ...props }) => {
          // Conflicts with List type
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
          const { type, ...rest } = props;
          return (
            <List component="ol" {...rest}>
              {children}
            </List>
          );
        },
        li: ({ children, ...props }) => <ListItem {...props}>{children}</ListItem>,
        table: ({ children, ...props }) => <Table {...props}>{children}</Table>,
        tbody: ({ children, ...props }) => <Tbody {...props}>{children}</Tbody>,
        thead: ({ children, ...props }) => <Thead {...props}>{children}</Thead>,
        tr: ({ children, ...props }) => <Tr {...props}>{children}</Tr>,
        td: ({ children, ...props }) => {
          // Conflicts with Td type
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
          const { width, ...rest } = props;
          return <Td {...rest}>{children}</Td>;
        },
        th: ({ children, ...props }) => <Th {...props}>{children}</Th>,
        img: ({ src, alt, ...props }) => {
          if (!src) {
            return null;
          }
          return <img src={src} alt={alt || 'Model documentation image'} {...props} />;
        },
      }}
      rehypePlugins={[rehypeRaw, rehypeUnwrapImages, rehypeSanitize]}
      remarkPlugins={[remarkGfm]}
    >
      {data}
    </Markdown>
  </div>
);

export default MarkdownComponent;
