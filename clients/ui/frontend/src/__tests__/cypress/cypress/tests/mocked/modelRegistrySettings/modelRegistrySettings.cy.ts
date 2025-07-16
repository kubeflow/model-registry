import type { ModelRegistryKind } from 'mod-arch-shared';
import { mockModelRegistryKind } from '~/__mocks__/mockModelRegistryKind';
import {
  modelRegistrySettings,
  FormFieldSelector,
} from '~/__tests__/cypress/cypress/pages/modelRegistrySettings';

const NAMESPACE = 'bella-namespace';
const roleBindingsMock = {
  data: {
    items: [
      {
        metadata: { name: 'stub-rb-1', creationTimestamp: null },
        roleRef: { apiGroup: '', kind: '', name: '' },
      },
      {
        metadata: { name: 'stub-rb-2', creationTimestamp: null },
        roleRef: { apiGroup: '', kind: '', name: '' },
      },
      {
        metadata: {
          name: 'model-registry-permissions',
          creationTimestamp: null,
          labels: {
            app: 'model-registry',
            'app.kubernetes.io/component': 'model-registry',
            'app.kubernetes.io/name': 'model-registry',
            'app.kubernetes.io/part-of': 'model-registry',
          },
        },
        subjects: [{ kind: 'User', apiGroup: 'rbac.authorization.k8s.io', name: 'admin-user' }],
        roleRef: {
          apiGroup: 'rbac.authorization.k8s.io',
          kind: 'Role',
          name: 'registry-user-model-registry',
        },
      },
      {
        metadata: {
          name: 'model-registry-dora-permissions',
          creationTimestamp: null,
          labels: {
            app: 'model-registry-dora',
            'app.kubernetes.io/component': 'model-registry',
            'app.kubernetes.io/name': 'model-registry-dora',
            'app.kubernetes.io/part-of': 'model-registry',
          },
        },
        subjects: [{ kind: 'User', apiGroup: 'rbac.authorization.k8s.io', name: 'dora-user' }],
        roleRef: {
          apiGroup: 'rbac.authorization.k8s.io',
          kind: 'Role',
          name: 'registry-user-model-registry-dora',
        },
      },
      {
        metadata: {
          name: 'model-registry-bella-permissions',
          creationTimestamp: null,
          labels: {
            app: 'model-registry-bella',
            'app.kubernetes.io/component': 'model-registry',
            'app.kubernetes.io/name': 'model-registry-bella',
            'app.kubernetes.io/part-of': 'model-registry',
          },
        },
        subjects: [{ kind: 'Group', apiGroup: 'rbac.authorization.k8s.io', name: 'bella-team' }],
        roleRef: {
          apiGroup: 'rbac.authorization.k8s.io',
          kind: 'Role',
          name: 'registry-user-model-registry-bella',
        },
      },
    ],
  },
};

const userMock = {
  data: {
    userId: 'user@example.com',
    clusterAdmin: true,
  },
};

const modelRegistryMock = mockModelRegistryKind({
  name: 'model-registry',
  displayName: 'Model Registry',
  description: 'Main registry',
});
const modelRegistryMockDora = mockModelRegistryKind({
  name: 'model-registry-dora',
  displayName: 'Dora Registry',
  description: 'Dora registry',
});
const modelRegistryMockBella = mockModelRegistryKind({
  name: 'model-registry-bella',
  displayName: 'Bella Registry',
  description: 'Bella registry',
});

const setupModelRegistryMocks = (
  registries: ModelRegistryKind[] = [
    modelRegistryMock,
    modelRegistryMockDora,
    modelRegistryMockBella,
  ],
) => {
  cy.intercept('GET', '/model-registry/api/v1/namespaces', {
    data: [{ metadata: { name: NAMESPACE } }],
  });
  cy.intercept('GET', '/model-registry/api/v1/user', userMock);
  cy.intercept('GET', '/model-registry/api/v1/settings/model_registry*', { data: registries });
  cy.intercept('GET', '/model-registry/api/v1/settings/role_bindings*', roleBindingsMock);
};

function selectNamespaceIfPresent() {
  cy.get('body').then(($body) => {
    if ($body.find('[data-testid="namespace-select"]').length) {
      cy.get('[data-testid="namespace-select"]').click();
      cy.findByText(NAMESPACE).click();
    }
  });
}

describe('Model Registry Settings', () => {
  it('should display the settings page', () => {
    setupModelRegistryMocks();
    modelRegistrySettings.visit();
    modelRegistrySettings.findHeading();
  });

  it('should show empty state when no registries exist', () => {
    setupModelRegistryMocks([]);
    modelRegistrySettings.visit();
    modelRegistrySettings.findEmptyState().should('exist');
  });

  it('should show table when there are registries', () => {
    setupModelRegistryMocks([
      mockModelRegistryKind({ name: 'test-registry-1' }),
      mockModelRegistryKind({ name: 'test-registry-2' }),
    ]);
    modelRegistrySettings.visit();
    modelRegistrySettings.findTable().should('exist');
    modelRegistrySettings.findModelRegistryRow('test-registry-1').should('exist');
  });

  describe('CreateModal', () => {
    beforeEach(() => {
      setupModelRegistryMocks([mockModelRegistryKind({ name: 'test-registry-1' })]);
    });
    it('should enable submit button if fields are valid', () => {
      modelRegistrySettings.visit();
      selectNamespaceIfPresent();
      modelRegistrySettings.findCreateButton().click({ force: true });
      modelRegistrySettings.findFormField(FormFieldSelector.NAME).type('valid-mr-name');
      modelRegistrySettings.findFormField(FormFieldSelector.HOST).type('localhost');
      modelRegistrySettings.findFormField(FormFieldSelector.PORT).type('5432');
      modelRegistrySettings.findFormField(FormFieldSelector.USERNAME).type('testuser');
      modelRegistrySettings.findFormField(FormFieldSelector.PASSWORD).type('testpass');
      modelRegistrySettings.findFormField(FormFieldSelector.DATABASE).type('testdb');
      modelRegistrySettings.shouldHaveNoErrors();
      modelRegistrySettings.findSubmitButton().should('be.enabled');
    });
  });

  //TODO: Add manage permission tests for model registry settings page

  describe('ManagePermissions', () => {
    it('should show the Manage permissions button for a model registry row', () => {
      setupModelRegistryMocks([mockModelRegistryKind({ name: 'model-registry' })]);
      modelRegistrySettings.visit();
      selectNamespaceIfPresent();
      modelRegistrySettings
        .findModelRegistryRow('model-registry')
        .findByText('Manage permissions')
        .should('be.visible');
    });
  });

  describe('DeleteModelRegistryModal', () => {
    beforeEach(() => {
      setupModelRegistryMocks([mockModelRegistryKind({ name: 'model-registry' })]);
      modelRegistrySettings.visit();
      selectNamespaceIfPresent();
      modelRegistrySettings
        .findModelRegistryRow('model-registry')
        .findKebabAction('Delete model registry')
        .click();
    });

    it('disables confirm button before name is typed', () => {
      cy.contains('button', 'Delete model registry').should('be.disabled');
    });

    it('enables confirm button after name is typed', () => {
      modelRegistrySettings.findConfirmDeleteNameInput().type('model-registry');
      modelRegistrySettings.findSubmitButton().should('be.enabled');
    });
  });
});
