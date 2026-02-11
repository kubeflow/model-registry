import { TempDevFeature } from '~/app/hooks/useTempDevFeatureAvailable';

class RegisterAndStoreFields {
  visit(enableRegistryStorageFeature = true, registryNamespace?: string) {
    if (enableRegistryStorageFeature) {
      window.localStorage.setItem(TempDevFeature.RegistryStorage, 'true');
    }
    const preferredModelRegistry = 'modelregistry-sample';
    const query = registryNamespace ? `?namespace=${encodeURIComponent(registryNamespace)}` : '';
    cy.visit(`/model-registry/${preferredModelRegistry}/register/model${query}`);
  }

  findNamespaceFormGroup() {
    return cy.findByTestId('namespace-form-group');
  }

  findNamespaceSelector() {
    return cy.findByTestId('form-namespace-selector');
  }

  /** Wrapper that contains the namespace Select trigger - use for disabled check */
  findNamespaceSelectTrigger() {
    return cy.findByTestId('form-namespace-selector-trigger');
  }

  /** MUI Select combobox inside form namespace selector - use to open/close dropdown */
  findNamespaceSelectCombobox() {
    return this.findNamespaceSelector().find('[role="combobox"]');
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
    this.findNamespaceSelectCombobox().scrollIntoView().click({ force: true });
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
    this.findNamespaceSelectCombobox().scrollIntoView().click({ force: true });
    namespaces.forEach((namespace) => {
      cy.findByRole('option', { name: namespace }).should('exist');
    });
    this.findNamespaceSelectCombobox().scrollIntoView().click({ force: true });
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

  findNamespaceRegistryAccessAlert() {
    return cy.findByTestId('namespace-registry-access-alert');
  }

  shouldShowNamespaceLabel() {
    this.findNamespaceFormGroup().find('label').should('contain.text', 'Namespace');
    return this;
  }

  shouldBeNamespaceSelectorDisabled() {
    this.findNamespaceSelectTrigger().find('[aria-disabled="true"]').should('exist');
    return this;
  }

  shouldShowNoAccessWarning() {
    this.findNamespaceRegistryAccessAlert()
      .should('be.visible')
      .and('contain.text', 'The selected namespace does not have access to this model registry');
    return this;
  }
}

export const registerAndStoreFields = new RegisterAndStoreFields();
