import React from 'react';
import { CodeBlock, CodeBlockCode } from '@patternfly/react-core';

type CodeBlockComponentProps = {
  children: React.ReactNode;
  className?: string;
};

const CodeBlockComponent: React.FC<CodeBlockComponentProps> = ({ children, className }) => (
  <CodeBlock className={className}>
    <CodeBlockCode>{children}</CodeBlockCode>
  </CodeBlock>
);

export default CodeBlockComponent;
