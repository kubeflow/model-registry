/* eslint-disable camelcase */
import type { Namespace } from 'mod-arch-core';
import { mockNamespace } from '~/__mocks__/mockNamespace';
import { appChrome } from '~/__tests__/cypress/cypress/pages/appChrome';
import { navBar } from '~/__tests__/cypress/cypress/pages/navBar';
import { mockUserSettings } from '~/__mocks__/mockUserSettings';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import type { ModelRegistry } from '~/app/types';

type HandlersProps = {
  modelRegistries?: ModelRegistry[];
  namespaces?: Namespace[];
  username?: string;
};

const initIntercepts = ({
  modelRegistries = [],
  username = 'test-user',
  namespaces = [
    mockNamespace({ name: 'namespace-1' }),
    mockNamespace({ name: 'namespace-2' }),
    mockNamespace({ name: 'namespace-3' }),
  ],
}: HandlersProps = {}) => {
  cy.interceptApi(
    'GET /api/:apiVersion/user',
    {
      path: {
        apiVersion: MODEL_REGISTRY_API_VERSION,
      },
    },
    mockUserSettings({
      userId: username,
    }),
  );

  cy.interceptApi(
    'GET /api/:apiVersion/namespaces',
    {
      path: {
        apiVersion: MODEL_REGISTRY_API_VERSION,
      },
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
};

describe('NavBar', () => {
  beforeEach(() => {
    cy.intercept('/logout').as('logout');
  });

  it('Should display empty state when no namespaces are returned', () => {
    initIntercepts({ namespaces: [] });
    appChrome.visit();
    navBar.shouldNamespaceSelectorHaveNoItems();
  });

  it('Should show username and log out', () => {
    initIntercepts({});
    appChrome.visit();
    navBar.findUsername().should('have.text', 'test-user');
  });

  it('Should select and update namespace', () => {
    initIntercepts({});
    appChrome.visit();

    navBar.findNamespaceSelector().findByText('namespace-1').should('exist');
    navBar.selectNamespace('namespace-2');
    navBar.findNamespaceSelector().findByText('namespace-2').should('exist');
  });
});

describe('NavBar - NamespaceSelector', () => {
  beforeEach(() => {
    cy.intercept('/logout').as('logout');
  });

  it('Should display empty state when no namespaces are returned', () => {
    initIntercepts({ namespaces: [] });
    appChrome.visit();
    navBar.shouldNamespaceSelectorHaveNoItems();
  });

  it('Should auto-select first namespace on initial load', () => {
    initIntercepts({
      namespaces: [
        mockNamespace({ name: 'namespace-1' }),
        mockNamespace({ name: 'namespace-2' }),
        mockNamespace({ name: 'namespace-3' }),
      ],
    });
    appChrome.visit();
    navBar.shouldNamespaceSelectorShow('namespace-1');
  });

  it('Should select and update namespace', () => {
    initIntercepts({});
    appChrome.visit();

    navBar.findNamespaceSelector().findByText('namespace-1').should('exist');
    navBar.selectNamespace('namespace-2');
    navBar.findNamespaceSelector().findByText('namespace-2').should('exist');
  });

  it('Should maintain namespace selection across navigation', () => {
    initIntercepts({});
    appChrome.visit();

    navBar.selectNamespace('namespace-2');
    navBar.shouldNamespaceSelectorShow('namespace-2');
    appChrome.findNavItem('Model Catalog').click();
    navBar.shouldNamespaceSelectorShow('namespace-2');
    appChrome.findNavItem('Model Registry').click();
    navBar.shouldNamespaceSelectorShow('namespace-2');
  });
});
