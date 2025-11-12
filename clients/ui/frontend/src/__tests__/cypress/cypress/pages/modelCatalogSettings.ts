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
}

export const modelCatalogSettings = new ModelCatalogSettings();
export const manageSourcePage = new ManageSourcePage();
