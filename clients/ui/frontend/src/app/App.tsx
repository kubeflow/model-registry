import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import AppRoutes from '@app/AppRoutes';
import '@app/app.css';
import {
  Flex,
  Masthead,
  MastheadContent,
  MastheadToggle,
  Page,
  PageToggleButton,
  Title
} from '@patternfly/react-core';
import NavSidebar from './NavSidebar';
import { BarsIcon } from '@patternfly/react-icons';

const App: React.FC = () => {
  const masthead = (
    <Masthead>
      <MastheadToggle>
        <PageToggleButton id="page-nav-toggle" variant="plain" aria-label="Dashboard navigation">
          <BarsIcon />
        </PageToggleButton>
      </MastheadToggle>

      <MastheadContent>
        <Flex>
          <Title headingLevel="h2" size="3xl">
            Kubeflow Model Registry UI
          </Title>
        </Flex>
      </MastheadContent>
    </Masthead>
  );

  return (
    <Page
      mainContainerId='primary-app-container'
      masthead={masthead}
      isManagedSidebar
      sidebar={<NavSidebar />}
    >
      <AppRoutes />
    </Page>
  );
};

export default App;
