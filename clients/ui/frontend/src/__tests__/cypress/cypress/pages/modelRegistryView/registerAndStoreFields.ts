import { TempDevFeature } from '~/app/hooks/useTempDevFeatureAvailable';

class RegisterAndStoreFields {
  visit(enableRegistryStorageFeature = true) {
    if (enableRegistryStorageFeature) {
      window.localStorage.setItem(TempDevFeature.RegistryStorage, 'true');
    }
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/register/model`);
  }

  findNamespaceFormGroup() {
    return cy.findByTestId('namespace-form-group');
  }

  findNamespaceSelector() {
    return cy.findByTestId('form-namespace-selector');
  }

  findOriginLocationSection() {
    return cy.findByTestId('model-origin-location-section');
  }

  findDestinationLocationSection() {
    return cy.findByTestId('model-destination-location-section');
  }

  findRegistrationModeToggleGroup() {
    return cy.findByTestId('registration-mode-toggle-group');
  }

  findRegisterToggleButton() {
    return cy.findByTestId('registration-mode-register');
  }

  findRegisterAndStoreToggleButton() {
    return cy.findByTestId('registration-mode-register-and-store');
  }

  selectNamespace(name: string) {
    this.findNamespaceSelector().click();
    cy.findByRole('option', { name }).click();
  }

  selectRegisterMode() {
    this.findRegisterToggleButton().click();
  }

  selectRegisterAndStoreMode() {
    this.findRegisterAndStoreToggleButton().click();
  }

  shouldShowPlaceholder(placeholder = 'Select a namespace') {
    this.findNamespaceSelector().findByText(placeholder).should('contain.text', placeholder);
    return this;
  }

  shouldHaveNamespaceOptions(namespaces: string[]) {
    this.findNamespaceSelector().click();
    namespaces.forEach((namespace) => {
      cy.findByRole('option', { name: namespace }).should('exist');
    });
    this.findNamespaceSelector().click();
    return this;
  }

  shouldShowSelectedNamespace(name: string) {
    this.findNamespaceSelector().findByText(name).should('have.text', name);
    return this;
  }

  shouldHideOriginLocationSection() {
    this.findOriginLocationSection().should('not.exist');
    return this;
  }

  shouldHideDestinationLocationSection() {
    this.findDestinationLocationSection().should('not.exist');
    return this;
  }

  shouldShowOriginLocationSection() {
    this.findOriginLocationSection().should('exist');
    return this;
  }

  shouldShowDestinationLocationSection() {
    this.findDestinationLocationSection().should('exist');
    return this;
  }

  shouldHaveRegistrationModeToggle() {
    this.findRegistrationModeToggleGroup().should('exist');
    return this;
  }

  shouldHaveRegisterModeSelected() {
    this.findRegisterToggleButton().find('button').should('have.attr', 'aria-pressed', 'true');
    return this;
  }

  shouldHaveRegisterAndStoreModeSelected() {
    this.findRegisterAndStoreToggleButton()
      .find('button')
      .should('have.attr', 'aria-pressed', 'true');
    return this;
  }
}

export const registerAndStoreFields = new RegisterAndStoreFields();
