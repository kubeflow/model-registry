import { appChrome } from './appChrome';

export enum FormFieldSelector {
  NAME = '#mr-name',
  RESOURCENAME = '#resource-mr-name',
  HOST = '#mr-host',
  PORT = '#mr-port',
  USERNAME = '#mr-username',
  PASSWORD = '#mr-password',
  DATABASE = '#mr-database',
}

export enum FormErrorTestId {
  HOST = 'mr-host-error',
  PORT = 'mr-port-error',
  USERNAME = 'mr-username-error',
  PASSWORD = 'mr-password-error',
  DATABASE = 'mr-database-error',
}

export enum DatabaseDetailsTestId {
  HOST = 'mr-db-host',
  PORT = 'mr-db-port',
  USERNAME = 'mr-db-username',
  PASSWORD = 'mr-db-password',
  DATABASE = 'mr-db-database',
}

class ModelRegistrySettings {
  visit(wait = true) {
    cy.visit('/model-registry-settings');
    if (wait) {
      this.wait();
    }
  }

  navigate() {
    this.findNavItem().click();
    this.wait();
  }

  private wait() {
    this.findHeading();
    cy.testA11y();
  }

  findHeading() {
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Model Registry Settings');
  }

  findNavItem() {
    return appChrome.findNavItem('Model registry settings', 'Settings');
  }

  findEmptyState() {
    return cy.findByTestId('mr-settings-empty-state');
  }

  findCreateButton() {
    return cy.findByText('Create model registry');
  }

  findFormField(selector: FormFieldSelector) {
    return cy.get(selector);
  }

  clearFormFields() {
    Object.values(FormFieldSelector).forEach((selector) => {
      this.findFormField(selector).clear();
      this.findFormField(selector).blur();
    });
  }

  findFormError(testId: FormErrorTestId) {
    return cy.findByTestId(testId);
  }

  shouldHaveAllErrors() {
    Object.values(FormErrorTestId).forEach((testId) => this.findFormError(testId).should('exist'));
  }

  shouldHaveNoErrors() {
    Object.values(FormErrorTestId).forEach((testId) =>
      this.findFormError(testId).should('not.exist'),
    );
  }

  findSubmitButton() {
    return cy.contains('button', 'Create');
  }

  findCancelButton() {
    return cy.findByTestId('modal-cancel-button');
  }

  findTable() {
    return cy.findByTestId('model-registries-table');
  }

  findModelRegistryRow(registryName: string) {
    return this.findTable().findByText(registryName).closest('tr');
  }

  findDatabaseDetail(testId: DatabaseDetailsTestId) {
    return cy.findByTestId(testId);
  }

  findDatabasePasswordHiddenButton() {
    return this.findDatabaseDetail(DatabaseDetailsTestId.PASSWORD).findByTestId(
      'password-hidden-button',
    );
  }

  findConfirmDeleteNameInput() {
    return cy.findByTestId('confirm-delete-input');
  }

  findManagePermissionsTooltip() {
    return cy.get('[data-testid="manage-permissions-tooltip"]');
  }
}

export const modelRegistrySettings = new ModelRegistrySettings();
