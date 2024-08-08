import * as React from 'react';
import { ExclamationTriangleIcon } from '@patternfly/react-icons';
import { Button, EmptyState, EmptyStateBody, EmptyStateFooter, PageSection } from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';

const NotFound: React.FunctionComponent = () => {
  function GoHomeBtn() {
    const navigate = useNavigate();

    function handleClick() {
      navigate('/');
    }

    return (
      <Button onClick={handleClick}>Take me home</Button>
    );
  }

  return (
    <PageSection>
      <EmptyState titleText="404 Page not found" variant="full" icon={ExclamationTriangleIcon}>
        <EmptyStateBody>
          We didn&apos;t find a page that matches the address you navigated to.
        </EmptyStateBody><EmptyStateFooter>
        <GoHomeBtn />
      </EmptyStateFooter></EmptyState>
    </PageSection>
  );
};

export { NotFound };
