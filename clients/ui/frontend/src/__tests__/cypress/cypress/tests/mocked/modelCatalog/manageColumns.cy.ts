/* eslint-disable camelcase */
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import {
  setupModelCatalogIntercepts,
  interceptPerformanceArtifactsList,
  interceptArtifactsList,
  type ModelCatalogInterceptOptions,
} from '~/__tests__/cypress/cypress/support/interceptHelpers/modelCatalog';
import { NBSP } from '~/__tests__/cypress/cypress/support/constants';

const initIntercepts = (options: Partial<ModelCatalogInterceptOptions> = {}) => {
  const resolvedOptions = {
    useValidatedModel: true,
    includePerformanceArtifacts: true,
    ...options,
  };

  setupModelCatalogIntercepts(resolvedOptions);
  interceptArtifactsList();
  interceptPerformanceArtifactsList();
};

const navigateToPerformanceTab = () => {
  modelCatalog.visit();
  modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
  modelCatalog.findModelCatalogDetailLink().first().click();
  modelCatalog.clickPerformanceInsightsTab();
  modelCatalog.findHardwareConfigurationTable().should('be.visible');
};

describe('Categorized Manage Columns Modal', () => {
  beforeEach(() => {
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    initIntercepts();
    navigateToPerformanceTab();
  });

  describe('Opening and Closing', () => {
    it('should open the modal when Customize columns button is clicked', () => {
      modelCatalog.findManageColumnsButton().should('be.visible').click();
      modelCatalog.findManageColumnsModal().should('be.visible');
    });

    it('should close the modal when Cancel button is clicked', () => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');

      modelCatalog.findManageColumnsCancelButton().click();
      modelCatalog.findManageColumnsModal().should('not.exist');
    });

    it('should close the modal when Update button is clicked', () => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');

      modelCatalog.findManageColumnsUpdateButton().click();
      modelCatalog.findManageColumnsModal().should('not.exist');
    });
  });

  describe('Modal Structure', () => {
    beforeEach(() => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');
    });

    it('should display the modal title and description', () => {
      modelCatalog.findManageColumnsModal().should('contain.text', 'Customize columns');
      modelCatalog
        .findManageColumnsModal()
        .should(
          'contain.text',
          'Manage the columns that appear in the hardware configuration table.',
        );
    });

    it('should display the search input', () => {
      modelCatalog.findManageColumnsSearch().should('be.visible');
    });

    it('should display selected count label', () => {
      modelCatalog.findManageColumnsSelectedCount().should('contain.text', 'selected');
    });

    it('should display Restore default columns button', () => {
      modelCatalog.findManageColumnsRestoreDefaults().should('be.visible');
    });

    it('should display Update and Cancel buttons', () => {
      modelCatalog.findManageColumnsUpdateButton().should('be.visible');
      modelCatalog.findManageColumnsCancelButton().should('be.visible');
    });
  });

  describe('Category Sections', () => {
    beforeEach(() => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');
    });

    it('should display all category sections', () => {
      modelCatalog.findManageColumnsSection('general').should('be.visible');
      modelCatalog.findManageColumnsSection('ttft-latency').should('be.visible');
      modelCatalog.findManageColumnsSection('e2e-latency').should('be.visible');
      modelCatalog.findManageColumnsSection('itl-latency').scrollIntoView().should('be.visible');
      modelCatalog.findManageColumnsSection('tps').scrollIntoView().should('be.visible');
    });

    it('should display category labels', () => {
      modelCatalog.findManageColumnsModal().should('contain.text', 'General');
      modelCatalog.findManageColumnsModal().should('contain.text', 'TTFT Latency');
      modelCatalog.findManageColumnsModal().should('contain.text', 'E2E Latency');
      modelCatalog.findManageColumnsModal().should('contain.text', 'ITL Latency');
      modelCatalog.findManageColumnsModal().should('contain.text', 'Throughput (TPS)');
    });

    it('should display columns within the General category', () => {
      modelCatalog.findManageColumnsSection('general').should('contain.text', 'Replicas');
      modelCatalog
        .findManageColumnsSection('general')
        .should('contain.text', `RPS${NBSP}per Replica`);
      modelCatalog.findManageColumnsSection('general').should('contain.text', 'Total RPS');
      modelCatalog
        .findManageColumnsSection('general')
        .should('contain.text', `Mean${NBSP}Input Tokens`);
      modelCatalog
        .findManageColumnsSection('general')
        .should('contain.text', `Mean${NBSP}Output Tokens`);
      modelCatalog.findManageColumnsSection('general').should('contain.text', 'vLLM Version');
    });
  });

  describe('Search Functionality', () => {
    beforeEach(() => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');
    });

    it('should filter columns by search term', () => {
      modelCatalog.findManageColumnsSearch().type('TTFT');

      // TTFT section should be visible
      modelCatalog.findManageColumnsSection('ttft-latency').should('be.visible');

      // Non-matching sections should be hidden
      modelCatalog.findManageColumnsSection('e2e-latency').should('not.exist');
      modelCatalog.findManageColumnsSection('itl-latency').should('not.exist');
    });

    it('should show empty state when no columns match the search', () => {
      modelCatalog.findManageColumnsSearch().type('nonexistent column xyz');

      modelCatalog.findManageColumnsModal().should('contain.text', 'No results found');
    });

    it('should show all categories again after clearing search', () => {
      modelCatalog.findManageColumnsSearch().type('TTFT');
      modelCatalog.findManageColumnsSection('e2e-latency').should('not.exist');

      // Clear the search
      modelCatalog.findManageColumnsSearch().find('button[aria-label="Reset"]').click();

      // All sections should reappear
      modelCatalog.findManageColumnsSection('general').should('be.visible');
      modelCatalog.findManageColumnsSection('ttft-latency').should('be.visible');
      modelCatalog.findManageColumnsSection('e2e-latency').should('be.visible');
    });

    it('should filter across categories when searching for a shared term', () => {
      modelCatalog.findManageColumnsSearch().type('Mean');

      // Sections with "Mean" columns should still be visible
      modelCatalog.findManageColumnsSection('ttft-latency').should('be.visible');
      modelCatalog.findManageColumnsSection('e2e-latency').should('be.visible');
      modelCatalog.findManageColumnsSection('itl-latency').should('be.visible');
      modelCatalog.findManageColumnsSection('tps').should('be.visible');
      modelCatalog.findManageColumnsSection('general').should('be.visible');
    });
  });

  describe('Column Toggle', () => {
    beforeEach(() => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');
    });

    it('should toggle a column checkbox off and on', () => {
      // Replicas is visible by default - uncheck it
      modelCatalog.findManageColumnsCheckbox('replicas').should('be.checked');
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsCheckbox('replicas').should('not.be.checked');

      // Check it again
      modelCatalog.findManageColumnsCheckbox('replicas').check();
      modelCatalog.findManageColumnsCheckbox('replicas').should('be.checked');
    });

    it('should update selected count when toggling columns', () => {
      // Get the initial count text
      modelCatalog
        .findManageColumnsSelectedCount()
        .invoke('text')
        .then((initialText) => {
          const initialCount = parseInt(initialText, 10);

          // Uncheck a currently checked column
          modelCatalog.findManageColumnsCheckbox('replicas').uncheck();

          modelCatalog
            .findManageColumnsSelectedCount()
            .should('contain.text', `${initialCount - 1}`);
        });
    });
  });

  describe('Restore Defaults', () => {
    beforeEach(() => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsModal().should('be.visible');
    });

    it('should become disabled after restoring defaults', () => {
      // The latency filter effect modifies columns from defaults on mount,
      // so restore defaults starts enabled. Click it to get to default state.
      modelCatalog.findManageColumnsRestoreDefaults().click();
      modelCatalog.findManageColumnsRestoreDefaults().should('be.disabled');
    });

    it('should be enabled after modifying column visibility from default state', () => {
      // First restore to defaults
      modelCatalog.findManageColumnsRestoreDefaults().click();
      modelCatalog.findManageColumnsRestoreDefaults().should('be.disabled');

      // Uncheck a default-visible column
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsRestoreDefaults().should('not.be.disabled');
    });

    it('should restore default columns when clicked', () => {
      // Uncheck a column that is currently checked
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsCheckbox('replicas').should('not.be.checked');

      // Click restore defaults
      modelCatalog.findManageColumnsRestoreDefaults().click();

      // Column should be checked again and button should be disabled
      modelCatalog.findManageColumnsCheckbox('replicas').should('be.checked');
      modelCatalog.findManageColumnsRestoreDefaults().should('be.disabled');
    });
  });

  describe('Apply Changes to Table', () => {
    it('should hide a column from the table after unchecking and updating', () => {
      // Replicas is always visible in the current column set
      modelCatalog.findHardwareConfigurationTableHeaders().should('contain.text', 'Replicas');

      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsUpdateButton().click();

      modelCatalog.findHardwareConfigurationTableHeaders().should('not.contain.text', 'Replicas');
    });

    it('should show a previously hidden column after checking and updating', () => {
      // First hide Replicas
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsUpdateButton().click();
      modelCatalog.findHardwareConfigurationTableHeaders().should('not.contain.text', 'Replicas');

      // Now show it again
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsCheckbox('replicas').check();
      modelCatalog.findManageColumnsUpdateButton().click();
      modelCatalog.findHardwareConfigurationTableHeaders().should('contain.text', 'Replicas');
    });

    it('should not apply changes when Cancel is clicked', () => {
      // Replicas column is visible
      modelCatalog.findHardwareConfigurationTableHeaders().should('contain.text', 'Replicas');

      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsCancelButton().click();

      // Column should still be visible since we cancelled
      modelCatalog.findHardwareConfigurationTableHeaders().should('contain.text', 'Replicas');
    });

    it('should show a non-default column in the table after enabling it', () => {
      // E2E Mean is not visible by default
      modelCatalog
        .findHardwareConfigurationTableHeaders()
        .should('not.contain.text', `E2E${NBSP}Latency Mean`);

      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsCheckbox('e2e_mean').check();
      modelCatalog.findManageColumnsUpdateButton().click();

      // E2E Mean should now be visible
      modelCatalog
        .findHardwareConfigurationTableHeaders()
        .should('contain.text', `E2E${NBSP}Latency Mean`);
    });
  });

  describe('Modal State Reset', () => {
    it('should reset unsaved changes when modal is reopened', () => {
      modelCatalog.findManageColumnsButton().click();

      // Uncheck replicas but cancel
      modelCatalog.findManageColumnsCheckbox('replicas').uncheck();
      modelCatalog.findManageColumnsCheckbox('replicas').should('not.be.checked');
      modelCatalog.findManageColumnsCancelButton().click();

      // Reopen modal - replicas should be checked again
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsCheckbox('replicas').should('be.checked');
    });

    it('should clear search when modal is reopened', () => {
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsSearch().type('TTFT');
      modelCatalog.findManageColumnsSection('e2e-latency').should('not.exist');
      modelCatalog.findManageColumnsCancelButton().click();

      // Reopen modal - search should be cleared and all sections visible
      modelCatalog.findManageColumnsButton().click();
      modelCatalog.findManageColumnsSearch().should('have.value', '');
      modelCatalog.findManageColumnsSection('e2e-latency').should('be.visible');
    });
  });
});
