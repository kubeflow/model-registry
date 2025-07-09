import type { ModelRegistryKind } from 'mod-arch-shared';
import { mockModelRegistryKind } from '~/__mocks__/mockModelRegistryKind';
import { modelRegistrySettings } from '~/__tests__/cypress/cypress/pages/modelRegistrySettings';

type HandlersProps = {
  modelRegistries?: ModelRegistryKind[];
};

const MODEL_REGISTRY_API_VERSION = 'v1';

const initIntercepts = ({
  modelRegistries = [
    mockModelRegistryKind({
      name: 'modelregistry-sample',
      description: 'New model registry',
      displayName: 'Model Registry Sample',
    }),
    mockModelRegistryKind({
      name: 'modelregistry-sample-2',
      description: 'New model registry 2',
      displayName: 'Model Registry Sample 2',
    }),
  ],
}: HandlersProps) => {
  cy.interceptApi(
    `GET /api/:apiVersion/settings/model_registry`,
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    modelRegistries.map((mr) => ({
      ...mr,
      metadata: {
        ...mr.metadata,
        annotations: {
          ...mr.metadata.annotations,
          'openshift.io/display-name': mr.metadata.displayName,
        },
      },
    })),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/settings/role_bindings`,
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    { items: [] },
  );
};

it('Shows empty state when there are no registries', () => {
  initIntercepts({ modelRegistries: [] });
  modelRegistrySettings.visit(true);
  modelRegistrySettings.findEmptyState().should('exist');
});

describe('ModelRegistriesTable', () => {
  it('Shows table when there are registries', () => {
    initIntercepts({});
    modelRegistrySettings.visit(true);
    modelRegistrySettings.findEmptyState().should('not.exist');
    modelRegistrySettings.findTable().should('exist');
    modelRegistrySettings.findModelRegistryRow('Model Registry Sample').should('exist');
  });
});
