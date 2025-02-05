import React from 'react';
import {
  Content,
  Divider,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import useModelRegistries from '~/app/hooks/useModelRegistries';
import useQueryParamNamespaces from '~/shared/hooks/useQueryParamNamespaces';
import ModelRegistriesTable from './ModelRegistriesTable';

const ModelRegistrySettings: React.FC = () => {
  const queryParams = useQueryParamNamespaces();
  const [modelRegistries, mrloaded, loadError] = useModelRegistries(queryParams);
  // TODO: [Midstream] Implement this when adding logic for rules review
  const loaded = mrloaded; //&& roleBindings.loaded;

  return (
    <>
      <ApplicationsPage
        title="Model Registry Settings"
        description={
          <Flex>
            <FlexItem>
              <Content>Manage model registry settings for all users in your organization.</Content>
            </FlexItem>
            <Divider />
          </Flex>
        }
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
            <EmptyStateBody>
              To get started, create a model registry. You can manage permissions after creation.
            </EmptyStateBody>
          </EmptyState>
        }
        provideChildrenPadding
      >
        <ModelRegistriesTable modelRegistries={modelRegistries} />
      </ApplicationsPage>
    </>
  );
};

export default ModelRegistrySettings;
