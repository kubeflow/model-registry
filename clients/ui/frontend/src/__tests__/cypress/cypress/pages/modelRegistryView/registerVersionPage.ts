export enum FormFieldSelector {
  REGISTERED_MODEL = '#registered-model-container .pf-m-typeahead',
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

class RegisterVersionPage {
  visit(registeredModelId?: string) {
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(
      registeredModelId
        ? `/model-registry/${preferredModelRegistry}/registeredModels/${registeredModelId}/registerVersion`
        : `/model-registry/${preferredModelRegistry}/registerVersion`,
    );
    this.wait();
  }

  private wait() {
    const preferredModelRegistry = 'modelregistry-sample';
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Register new version');
    cy.findByText(`Model registry - ${preferredModelRegistry}`).should('exist');
    cy.testA11y();
  }

  findFormField(selector: FormFieldSelector) {
    return cy.get(selector);
  }

  selectRegisteredModel(name: string) {
    this.findFormField(FormFieldSelector.REGISTERED_MODEL)
      .findByRole('button', { name: 'Typeahead menu toggle' })
      .findSelectOption(name)
      .click();
  }

  findSubmitButton() {
    return cy.findByTestId('create-button');
  }
}

export const registerVersionPage = new RegisterVersionPage();
