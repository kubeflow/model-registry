import { SearchSelector } from '~/__tests__/cypress/cypress/pages/components/subComponents/SearchSelector';

export enum FormFieldSelector {
  MODEL_NAME = '#model-name',
  MODEL_DESCRIPTION = '#model-description',
  VERSION_NAME = '#version-name',
  VERSION_DESCRIPTION = '#version-description',
  SOURCE_MODEL_FORMAT = '#source-model-format',
  SOURCE_MODEL_FORMAT_VERSION = '#source-model-format-version',
  LOCATION_TYPE_OBJECT_STORAGE = '#location-type-object-storage',
  LOCATION_ENDPOINT = '#location-endpoint',
  LOCATION_BUCKET = '#location-bucket',
  LOCATION_REGION = '#location-region',
  LOCATION_PATH = '#location-path',
  LOCATION_TYPE_URI = '#location-type-uri',
  LOCATION_URI = '#location-uri',
}

class RegisterModelPage {
  projectDropdown = new SearchSelector('project-selector', 'connection-autofill-modal');

  visit() {
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registerModel`);
    this.wait();
  }

  private wait() {
    const preferredModelRegistry = 'modelregistry-sample';
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Register model');
    cy.findByText(`Model registry - ${preferredModelRegistry}`).should('exist');
    cy.testA11y();
  }

  findFormField(selector: FormFieldSelector) {
    return cy.get(selector);
  }

  findObjectStorageAutofillButton() {
    return cy.findByTestId('object-storage-autofill-button');
  }

  findConnectionAutofillModal() {
    return cy.findByTestId('connection-autofill-modal');
  }

  findConnectionSelector() {
    return this.findConnectionAutofillModal().findByTestId('select-data-connection');
  }

  findAutofillButton() {
    return cy.findByTestId('autofill-modal-button');
  }

  findSubmitButton() {
    return cy.findByTestId('create-button');
  }
}

export const registerModelPage = new RegisterModelPage();
