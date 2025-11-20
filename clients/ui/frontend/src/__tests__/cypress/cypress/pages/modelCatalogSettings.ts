import { appChrome } from './appChrome';
import { TableRow } from './components/table';

class CatalogSourceConfigRow extends TableRow {
  findName() {
    return this.find().find('[data-label="Name"]');
  }

  findOrganization() {
    return this.find().find('[data-label="Organization"]');
  }

  findModelVisibility() {
    return this.find().find('[data-label="Model visibility"]');
  }

  findSourceType() {
    return this.find().find('[data-label="Source type"]');
  }

  findEnableToggle() {
    return this.find().find('[data-label="Enable"]').find('input[type="checkbox"]');
  }

  findValidationStatus() {
    return this.find().find('[data-label="Validation status"]');
  }

  findManageSourceButton() {
    return this.find()
      .find('[data-label="Actions"]')
      .findByRole('button', { name: 'Manage source' });
  }

  shouldHaveModelVisibility(visibility: 'Filtered' | 'Unfiltered') {
    this.findModelVisibility().contains(visibility);
    return this;
  }

  shouldHaveOrganization(org: string) {
    this.findOrganization().contains(org);
    return this;
  }

  shouldHaveSourceType(type: string) {
    this.findSourceType().contains(type);
    return this;
  }

  toggleEnable() {
    this.findEnableToggle().click();
    return this;
  }

  shouldHaveEnableToggle(shouldExist: boolean) {
    if (shouldExist) {
      this.findEnableToggle().should('exist');
    } else {
      this.find().find('[data-label="Enable"]').should('be.empty');
    }
    return this;
  }

  shouldHaveEnableState(enabled: boolean) {
    if (enabled) {
      this.findEnableToggle().should('be.checked');
    } else {
      this.findEnableToggle().should('not.be.checked');
    }
    return this;
  }
}

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

  findTable() {
    return cy.findByTestId('catalog-source-configs-table');
  }

  findEmptyState() {
    return cy.findByTestId('catalog-settings-empty-state');
  }

  getRow(name: string) {
    return new CatalogSourceConfigRow(() =>
      this.findTable().find('tbody').find('tr').contains(name).parents('tr'),
    );
  }

  findRows() {
    return this.findTable().find('tbody tr');
  }

  shouldHaveSourceConfigs() {
    this.findTable().should('exist');
    this.findRows().should('have.length.at.least', 1);
    return this;
  }

  shouldBeEmpty() {
    this.findEmptyState().should('exist');
    return this;
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
