import React from 'react';
import { EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import ApplicationsPage from '~/app/components/ApplicationsPage';

const ModelRegistrySettings: React.FC = () => {
  const [modelRegistries, loaded, loadError] = [[], true, undefined]; // TODO: change to real values
  return (
    <>
      <ApplicationsPage
        title="Model registry settings"
        description="Manage model registry settings for all users in your organization."
        loaded={loaded}
        loadError={loadError}
        errorMessage="Unable to load model registries."
        empty={modelRegistries.length === 0}
        emptyStatePage={
          <EmptyState
            headingLevel="h5"
            icon={PlusCircleIcon}
            titleText="No model registries"
            variant={EmptyStateVariant.lg}
            data-testid="mr-settings-empty-state"
          >
            <EmptyStateBody>To get started, create a model registry.</EmptyStateBody>
          </EmptyState>
        }
        provideChildrenPadding
      >
        TODO: Add model registry settings
      </ApplicationsPage>
    </>
  );
};

export default ModelRegistrySettings;
