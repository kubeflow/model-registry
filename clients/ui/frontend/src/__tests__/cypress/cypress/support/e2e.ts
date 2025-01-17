// ***********************************************************
// This example support/e2e.ts is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

import chaiSubset from 'chai-subset';
import '@cypress/code-coverage/support';
import { mockUserSettings } from '~/__mocks__/mockUserSettings';
import 'cypress-mochawesome-reporter/register';
import './commands';
import { mockNamespace } from '~/__mocks__/mockNamespace';
import { MODEL_REGISTRY_API_VERSION } from './commands/api';

chai.use(chaiSubset);

Cypress.Keyboard.defaults({
  keystrokeDelay: 0,
});

beforeEach(() => {
  if (Cypress.env('MOCK')) {
    // fallback: return 404 for all api requests
    cy.intercept({ pathname: '/api/**' }, { statusCode: 404 });

    cy.interceptApi(
      'GET /api/:apiVersion/user',
      {
        path: {
          apiVersion: MODEL_REGISTRY_API_VERSION,
        },
      },
      mockUserSettings({}),
    );

    cy.interceptApi(
      'GET /api/:apiVersion/namespaces',
      {
        path: {
          apiVersion: MODEL_REGISTRY_API_VERSION,
        },
      },
      [mockNamespace({})],
    );
  }
});
