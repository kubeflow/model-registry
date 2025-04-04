import React from 'react';
import { Divider, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import ApplicationsPage from '~/shared/components/ApplicationsPage';
import useModelRegistries from '~/app/hooks/useModelRegistries';
import useQueryParamNamespaces from '~/shared/hooks/useQueryParamNamespaces';
import { isMUITheme } from '~/shared/utilities/const';
import TitleWithIcon from '~/shared/components/design/TitleWithIcon';
import { ProjectObjectType } from '~/shared/components/design/utils';
// import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import ModelRegistriesTable from './ModelRegistriesTable';
import CreateModal from './ModelRegistryCreateModal';

const ModelRegistrySettings: React.FC = () => {
  const queryParams = useQueryParamNamespaces();
  const [
    modelRegistries,
    mrloaded,
    loadError,
    // refreshModelRegistries
  ] = useModelRegistries(queryParams);
  const [createModalOpen, setCreateModalOpen] = React.useState(false);
  // TODO: [Midstream] Implement this when adding logic for rules review
  // const { refreshRulesReview } = React.useContext(ModelRegistrySelectorContext);

  const loaded = mrloaded; //&& roleBindings.loaded;

  // TODO: implement when refreshModelRegistries() and refreshRulesReview() are added
  // const refreshAll = React.useCallback(
  //   () => Promise.all([refreshModelRegistries(), refreshRulesReview()]),
  //   [refreshModelRegistries, refreshRulesReview],
  // );

  return (
    <>
      <ApplicationsPage
        title={
          !isMUITheme() ? (
            <TitleWithIcon
              title="Model Registry Settings"
              objectType={ProjectObjectType.modelRegistrySettings}
            />
          ) : (
            'Model Registry Settings'
          )
        }
        description={
          !isMUITheme() ? (
            'Manage model registry settings for all users in your organization.'
          ) : (
            <Divider />
          )
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
        <ModelRegistriesTable
          modelRegistries={modelRegistries}
          onCreateModelRegistryClick={() => {
            setCreateModalOpen(true);
          }}
          // eslint-disable-next-line @typescript-eslint/no-empty-function
          refresh={() => {}}
        />
      </ApplicationsPage>
      {createModalOpen ? (
        <CreateModal
          onClose={() => setCreateModalOpen(false)}
          // refresh={refreshAll}
        />
      ) : null}
    </>
  );
};

export default ModelRegistrySettings;
