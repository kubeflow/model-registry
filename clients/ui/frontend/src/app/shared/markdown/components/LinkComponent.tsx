import React from 'react';
import { Button } from '@patternfly/react-core';

type LinkComponentProps = {
  children: React.ReactNode;
  href?: string;
  className?: string;
};

const LinkComponent: React.FC<LinkComponentProps> = ({ children, href = '#', className }) => (
  <Button
    variant="link"
    className={className}
    isInline
    onClick={() => href !== '#' && window.open(href, '_blank', 'noopener,noreferrer')}
  >
    {children}
  </Button>
);

export default LinkComponent;
