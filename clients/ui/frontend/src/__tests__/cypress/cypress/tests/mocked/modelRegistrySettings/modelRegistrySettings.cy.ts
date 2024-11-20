import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import type { ModelRegistry } from '~/app/types';
import { modelRegistrySettings } from '~/__tests__/cypress/cypress/pages/modelRegistrySettings';

type HandlersProps = {
  modelRegistries?: ModelRegistry[];
};

const MODEL_REGISTRY_API_VERSION = 'v1';

const initIntercepts = ({
  modelRegistries = [
    mockModelRegistry({
      name: 'modelregistry-sample',
      description: 'New model registry',
      displayName: 'Model Registry Sample',
    }),
    mockModelRegistry({
      name: 'modelregistry-sample-2',
      description: 'New model registry 2',
      displayName: 'Model Registry Sample 2',
    }),
  ],
}: HandlersProps) => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_registry`,
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    modelRegistries,
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
