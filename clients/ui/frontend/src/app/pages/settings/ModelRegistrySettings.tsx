import React from 'react';
import { EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import useModelRegistries from '~/app/hooks/useModelRegistries';
import TitleWithIcon from '~/shared/components/design/TitleWithIcon';
import { ProjectObjectType } from '~/shared/components/design/utils';
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
        title={
          <TitleWithIcon
            title="Model Registry Settings"
            objectType={ProjectObjectType.modelRegistrySettings}
          />
        }
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
