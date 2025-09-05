import React, { useState } from 'react';
import { Button, Content, CodeBlock, CodeBlockCode } from '@patternfly/react-core';
import { AngleDownIcon, AngleRightIcon } from '@patternfly/react-icons';

type DetailsComponentProps = {
  children: React.ReactNode;
  summary: string;
  className?: string;
};

const DetailsComponent: React.FC<DetailsComponentProps> = ({ children, summary, className }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const processedSummary = summary.replace(/\\n/g, ' ').trim();
  const processedChildren = React.Children.map(children, (child) => {
    if (React.isValidElement(child) && child.type === 'code') {
      return (
        <CodeBlock>
          <CodeBlockCode>{String(child.props.children).replace(/\\n/g, '\n')}</CodeBlockCode>
        </CodeBlock>
      );
    }
    return child;
  });

  return (
    <div className={className}>
      <Button
        variant="link"
        onClick={() => setIsExpanded(!isExpanded)}
        style={{ padding: 0, margin: 0, textAlign: 'left' }}
      >
        {isExpanded ? <AngleDownIcon /> : <AngleRightIcon />} {processedSummary}
      </Button>
      {isExpanded && (
        <Content style={{ marginLeft: '1.5rem', marginTop: '0.5rem' }}>{processedChildren}</Content>
      )}
    </div>
  );
};

export default DetailsComponent;
