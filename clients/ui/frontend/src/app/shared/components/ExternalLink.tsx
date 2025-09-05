import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';

type ExternalLinkProps = {
  text: string;
  to: string;
  testId?: string;
};

const ExternalLink: React.FC<ExternalLinkProps> = ({ text, to, testId }) => (
  <Button
    variant="link"
    data-testid={testId}
    isInline
    onClick={() => {
      window.open(to);
    }}
    icon={<ExternalLinkAltIcon />}
    iconPosition="end"
  >
    {text}
  </Button>
);

export default ExternalLink;
