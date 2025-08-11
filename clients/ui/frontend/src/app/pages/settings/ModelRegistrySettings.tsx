import React from 'react';
import {
  Divider,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import { ProjectObjectType, TitleWithIcon, ApplicationsPage } from 'mod-arch-shared';
import { useQueryParamNamespaces } from 'mod-arch-core';
// import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import useModelRegistriesSettings from '~/app/hooks/useModelRegistriesSetting';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import ModelRegistriesTable from './ModelRegistriesTable';
import CreateModal from './CreateModal';

const ModelRegistrySettings: React.FC = () => {
  const queryParams = useQueryParamNamespaces();
  const [
    modelRegistries,
    mrloaded,
    loadError,
    // refreshModelRegistries
  ] = useModelRegistriesSettings(queryParams);
  const roleBindings = useModelRegistryRoleBindings(queryParams);
  const [createModalOpen, setCreateModalOpen] = React.useState(false);
  // TODO: [Midstream] Implement this when adding logic for rules review
  // const { refreshRulesReview } = React.useContext(ModelRegistrySelectorContext);

  const loaded = mrloaded && roleBindings.loaded;

  // TODO: implement when refreshModelRegistries() and refreshRulesReview() are added
  // const refreshAll = React.useCallback(
  //   () => Promise.all([refreshModelRegistries(), refreshRulesReview()]),
  //   [refreshModelRegistries, refreshRulesReview],
  // );

  return (
    <>
      <ApplicationsPage
        title={
          <TitleWithIcon
            title="Model Registry Settings"
            objectType={ProjectObjectType.modelRegistrySettings}
          />
        }
        description={
          <Stack hasGutter>
            <StackItem>
              Manage model registry settings for all users in your organization.
            </StackItem>
            <StackItem>
              <Divider />
            </StackItem>
          </Stack>
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
          roleBindings={roleBindings}
          onCreateModelRegistryClick={() => {
            setCreateModalOpen(true);
          }}
          // eslint-disable-next-line @typescript-eslint/no-empty-function
          refresh={() => Promise.resolve()}
        />
      </ApplicationsPage>
      {createModalOpen ? (
        // eslint-disable-next-line @typescript-eslint/no-empty-function
        <CreateModal onClose={() => setCreateModalOpen(false)} refresh={() => Promise.resolve()} />
      ) : null}
    </>
  );
};

export default ModelRegistrySettings;
