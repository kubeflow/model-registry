import { TempDevFeature } from '~/app/hooks/useTempDevFeatureAvailable';
import { FormFieldSelector } from './registerModelPage';

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

  // Destination field finders (using id selectors since these inputs use id, not data-testid)
  findDestinationOciRegistryInput() {
    return cy.get('#destination-oci-registry');
  }

  findDestinationOciUriInput() {
    return cy.get('#destination-oci-uri');
  }

  findDestinationOciUsernameInput() {
    return cy.get('#destination-oci-username');
  }

  findDestinationOciPasswordInput() {
    return cy.get('#destination-oci-password');
  }

  // Source credential field finders
  findSourceS3AccessKeyIdInput() {
    return cy.get('#location-s3-access-key-id');
  }

  findSourceS3SecretAccessKeyInput() {
    return cy.get('#location-s3-secret-access-key');
  }

  // Submit button
  findSubmitButton() {
    return cy.findByTestId('create-button');
  }

  // Fill form helper methods
  fillModelName(name: string) {
    cy.get(FormFieldSelector.MODEL_NAME).clear();
    cy.get(FormFieldSelector.MODEL_NAME).type(name);
  }

  fillVersionName(name: string) {
    cy.get(FormFieldSelector.VERSION_NAME).clear();
    cy.get(FormFieldSelector.VERSION_NAME).type(name);
  }

  fillJobName(name: string) {
    cy.get(FormFieldSelector.JOB_NAME).clear();
    cy.get(FormFieldSelector.JOB_NAME).type(name);
  }

  fillSourceEndpoint(endpoint: string) {
    cy.get(FormFieldSelector.LOCATION_ENDPOINT).clear();
    cy.get(FormFieldSelector.LOCATION_ENDPOINT).type(endpoint);
  }

  fillSourceBucket(bucket: string) {
    cy.get(FormFieldSelector.LOCATION_BUCKET).clear();
    cy.get(FormFieldSelector.LOCATION_BUCKET).type(bucket);
  }

  fillSourcePath(path: string) {
    cy.get(FormFieldSelector.LOCATION_PATH).clear();
    cy.get(FormFieldSelector.LOCATION_PATH).type(path);
  }

  fillDestinationOciRegistry(registry: string) {
    this.findDestinationOciRegistryInput().clear();
    this.findDestinationOciRegistryInput().type(registry);
  }

  fillDestinationOciUri(uri: string) {
    this.findDestinationOciUriInput().clear();
    this.findDestinationOciUriInput().type(uri);
  }

  fillDestinationOciUsername(username: string) {
    this.findDestinationOciUsernameInput().clear();
    this.findDestinationOciUsernameInput().type(username);
  }

  fillDestinationOciPassword(password: string) {
    this.findDestinationOciPasswordInput().clear();
    this.findDestinationOciPasswordInput().type(password);
  }

  fillSourceS3AccessKeyId(accessKeyId: string) {
    this.findSourceS3AccessKeyIdInput().clear();
    this.findSourceS3AccessKeyIdInput().type(accessKeyId);
  }

  fillSourceS3SecretAccessKey(secretAccessKey: string) {
    this.findSourceS3SecretAccessKeyInput().clear();
    this.findSourceS3SecretAccessKeyInput().type(secretAccessKey);
  }

  // Convenience method to fill all required fields for submission
  fillAllRequiredFields() {
    this.fillModelName('test-model');
    this.fillVersionName('v1.0.0');
    this.fillJobName('my-transfer-job');
    this.fillSourceEndpoint('https://s3.amazonaws.com');
    this.fillSourceBucket('test-bucket');
    this.fillSourcePath('models/test');
    this.fillSourceS3AccessKeyId('AKIAIOSFODNN7EXAMPLE');
    this.fillSourceS3SecretAccessKey('wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY');
    this.fillDestinationOciRegistry('quay.io');
    this.fillDestinationOciUri('quay.io/my-org/my-model:v1');
    this.fillDestinationOciUsername('testuser');
    this.fillDestinationOciPassword('testpassword123');
  }
}

export const registerAndStoreFields = new RegisterAndStoreFields();
