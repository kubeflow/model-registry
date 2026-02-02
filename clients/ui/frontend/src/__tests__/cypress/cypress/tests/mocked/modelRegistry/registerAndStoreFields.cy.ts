/* eslint-disable camelcase */
import type { Namespace } from 'mod-arch-core';
import { mockNamespace } from '~/__mocks__/mockNamespace';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { registerAndStoreFields } from '~/__tests__/cypress/cypress/pages/modelRegistryView/registerAndStoreFields';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import type { ModelRegistry, RegisteredModel } from '~/app/types';
import { mockRegisteredModelList } from '~/__mocks__/mockRegisteredModelsList';

type HandlersProps = {
  modelRegistries?: ModelRegistry[];
  namespaces?: Namespace[];
  registeredModels?: RegisteredModel[];
};

const initIntercepts = ({
  modelRegistries = [mockModelRegistry({ name: 'modelregistry-sample' })],
  namespaces = [
    mockNamespace({ name: 'namespace-1' }),
    mockNamespace({ name: 'namespace-2' }),
    mockNamespace({ name: 'namespace-3' }),
  ],
  registeredModels = [],
}: HandlersProps = {}) => {
  cy.interceptApi(
    'GET /api/:apiVersion/namespaces',
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    namespaces,
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry`,
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    modelRegistries,
  );

  cy.interceptApi(
    'GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models',
    {
      path: {
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelRegistryName: 'modelregistry-sample',
      },
    },
    mockRegisteredModelList({ items: registeredModels }),
  );
};

describe('Register and Store Fields - Toggle Behavior', () => {
  beforeEach(() => {
    initIntercepts({});
    registerAndStoreFields.visit();
  });

  it('Should display registration mode toggle when feature is enabled', () => {
    registerAndStoreFields.shouldHaveRegistrationModeToggle();
  });

  it('Should have "Register" mode selected by default', () => {
    registerAndStoreFields.shouldHaveRegisterModeSelected();
  });

  it('Should switch to "Register and store" mode', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.shouldHaveRegisterAndStoreModeSelected();
  });

  it('Should show namespace selector only in "Register and store" mode', () => {
    registerAndStoreFields.findNamespaceFormGroup().should('not.exist');
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.findNamespaceSelector().should('exist');
  });

  it('Should switch back to "Register" mode and hide namespace selector', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.findNamespaceSelector().should('exist');
    registerAndStoreFields.selectRegisterMode();
    registerAndStoreFields.shouldHaveRegisterModeSelected();
    registerAndStoreFields.findNamespaceSelector().should('not.exist');
  });

  it('Should reset namespace selection when switching modes', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    registerAndStoreFields.shouldShowSelectedNamespace('namespace-1');
    registerAndStoreFields.selectRegisterMode();
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.shouldShowPlaceholder('Select a namespace');
  });
});

describe('Register and Store Fields - NamespaceSelector', () => {
  beforeEach(() => {
    initIntercepts({});
    registerAndStoreFields.visit();
    registerAndStoreFields.selectRegisterAndStoreMode();
  });

  it('Should show placeholder text instead of auto-selecting', () => {
    registerAndStoreFields.shouldShowPlaceholder('Select a namespace');
  });

  it('Should display all available namespaces in dropdown', () => {
    registerAndStoreFields.shouldHaveNamespaceOptions([
      'namespace-1',
      'namespace-2',
      'namespace-3',
    ]);
  });

  it('Should hide form sections until namespace is selected', () => {
    registerAndStoreFields.shouldHideOriginLocationSection().shouldHideDestinationLocationSection();
  });

  it('Should show form sections after namespace selection', () => {
    registerAndStoreFields.selectNamespace('namespace-1');

    registerAndStoreFields.shouldShowOriginLocationSection();
    registerAndStoreFields.shouldShowDestinationLocationSection();
  });

  it('Should update selected namespace in dropdown', () => {
    registerAndStoreFields.selectNamespace('namespace-2');
    registerAndStoreFields.shouldShowSelectedNamespace('namespace-2');
  });

  it('Should handle empty namespace list gracefully', () => {
    initIntercepts({ namespaces: [] });
    registerAndStoreFields.visit();
    registerAndStoreFields.selectRegisterAndStoreMode();

    registerAndStoreFields.findNamespaceSelector().should('exist');
    registerAndStoreFields.findNamespaceSelector().should('be.disabled');

    registerAndStoreFields.shouldShowPlaceholder('Select a namespace');
  });
});
