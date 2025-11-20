import { appChrome } from './appChrome';

class ModelCatalogSettings {
  visit(wait = true) {
    cy.visit('/model-catalog-settings');
    if (wait) {
      this.wait();
    }
  }

  navigate() {
    cy.get('body').then(($body) => {
      if ($body.find('#page-sidebar').length > 0) {
        this.findNavItem().click();
        this.wait();
      }
    });
  }

  private wait() {
    this.findHeading();
    cy.testA11y();
  }

  findHeading() {
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Model catalog settings');
  }

  findNavItem() {
    return appChrome.findNavItem('Model catalog settings', 'Settings');
  }

  findDescription() {
    return cy.contains('Manage model catalog sources for your organization.');
  }

  findAddSourceButton() {
    return cy.findByTestId('add-source-button');
  }
}

class ManageSourcePage {
  visitAddSource(wait = true) {
    cy.visit('/model-catalog-settings/add-source');
    if (wait) {
      this.wait();
    }
  }

  visitManageSource(catalogSourceId: string, wait = true) {
    cy.visit(`/model-catalog-settings/manage-source/${encodeURIComponent(catalogSourceId)}`);
    if (wait) {
      this.wait();
    }
  }

  private wait() {
    this.findHeading();
    cy.testA11y();
  }

  findHeading() {
    cy.findByTestId('app-page-title').should('exist');
  }

  findBreadcrumb() {
    return cy.get('a[href="/model-catalog-settings"]').contains('Model catalog settings');
  }

  findBreadcrumbAction() {
    return cy.findByTestId('breadcrumb-source-action');
  }

  findAddSourceTitle() {
    return cy.findByTestId('app-page-title').contains('Add a source');
  }

  findManageSourceTitle() {
    return cy.findByTestId('app-page-title').contains('Manage source');
  }

  findAddSourceDescription() {
    return cy.contains('Add a new model catalog source to your organization.');
  }

  findManageSourceDescription() {
    return cy.contains('Manage the selected model catalog source.');
  }

  // Form field methods
  findNameInput() {
    return cy.findByTestId('source-name-input');
  }

  findNameError() {
    return cy.findByTestId('source-name-error');
  }

  findSourceTypeHuggingFace() {
    return cy.findByTestId('source-type-huggingface');
  }

  findSourceTypeYaml() {
    return cy.findByTestId('source-type-yaml');
  }

  findSourceTypeHuggingFaceLabel() {
    return cy.get('label[for="source-type-huggingface"]');
  }

  findSourceTypeYamlLabel() {
    return cy.get('label[for="source-type-yaml"]');
  }

  findCredentialsSection() {
    return cy.findByTestId('credentials-section');
  }

  findAccessTokenInput() {
    return cy.findByTestId('access-token-input');
  }

  findAccessTokenError() {
    return cy.findByTestId('access-token-error');
  }

  findOrganizationInput() {
    return cy.findByTestId('organization-input');
  }

  findOrganizationError() {
    return cy.findByTestId('organization-error');
  }

  findYamlSection() {
    return cy.findByTestId('yaml-section');
  }

  findYamlContentInput() {
    return cy.findByTestId('yaml-content-input').find('textarea');
  }

  findYamlContentError() {
    return cy.findByTestId('yaml-content-error');
  }

  findModelVisibilitySection() {
    return cy.findByTestId('model-visibility-section');
  }

  toggleModelVisibility() {
    this.findModelVisibilitySection().find('button').first().click();
  }

  findAllowedModelsInput() {
    return cy.findByTestId('allowed-models-input');
  }

  findExcludedModelsInput() {
    return cy.findByTestId('excluded-models-input');
  }

  findEnableSourceCheckbox() {
    return cy.findByTestId('enable-source-checkbox');
  }

  findSubmitButton() {
    return cy.findByTestId('submit-button');
  }

  findPreviewButton() {
    return cy.findByTestId('preview-button');
  }

  findCancelButton() {
    return cy.findByTestId('cancel-button');
  }

  fillSourceName(name: string) {
    this.findNameInput().clear().type(name);
  }

  selectSourceType(type: 'huggingface' | 'yaml') {
    if (type === 'huggingface') {
      this.findSourceTypeHuggingFaceLabel().click();
    } else {
      this.findSourceTypeYamlLabel().click();
    }
  }

  fillAccessToken(token: string) {
    this.findAccessTokenInput().clear().type(token);
  }

  fillOrganization(org: string) {
    this.findOrganizationInput().clear().type(org);
  }

  fillYamlContent(yaml: string) {
    this.findYamlContentInput().clear().type(yaml);
  }

  fillAllowedModels(models: string) {
    this.findAllowedModelsInput().clear().type(models);
  }

  fillExcludedModels(models: string) {
    this.findExcludedModelsInput().clear().type(models);
  }

  toggleEnableSource() {
    this.findEnableSourceCheckbox().click();
  }

  findPreviewPanel() {
    return cy.findByTestId('preview-panel');
  }

  findPreviewPanelTitle() {
    return cy.contains('Model catalog preview');
  }

  findPreviewPanelEmptyMessage() {
    return cy.contains('To view the models from this source that will appear');
  }

  findPreviewButtonHeader() {
    return cy.findByTestId('preview-button-header');
  }

  findPreviewButtonPanel() {
    return cy.findByTestId('preview-button-panel');
  }
}

export const modelCatalogSettings = new ModelCatalogSettings();
export const manageSourcePage = new ManageSourcePage();
