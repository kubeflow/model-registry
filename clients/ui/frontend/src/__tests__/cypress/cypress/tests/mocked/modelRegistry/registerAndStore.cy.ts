import { mockModelRegistryList, mockRegisteredModelList } from '~/__mocks__';
import { mockDashboardConfig } from '~/__mocks__/mockDashboardConfig';
import { mockStatus } from '~/__mocks__/mockStatus';
import { TempDevFeature } from '~/app/hooks/useTempDevFeatureAvailable';
import {
  registerModelPage,
  FormFieldSelector as RegisterModelFormFieldSelector,
} from '~/__tests__/cypress/cypress/pages/modelRegistryView/registerModelPage';
import {
  registerVersionPage,
  FormFieldSelector as RegisterVersionFormFieldSelector,
} from '~/__tests__/cypress/cypress/pages/modelRegistryView/registerVersionPage';

const initIntercepts = () => {
  cy.interceptOdh(
    'GET /api/config',
    mockDashboardConfig({
      // Enable the RegistryStorage feature flag
      disabledCustomizationFeatures: [],
      enabledCustomizationFeatures: [TempDevFeature.RegistryStorage],
    }),
  );
  cy.interceptK8sList(
    { model: { apiVersion: 'v1alpha3', kind: 'ModelRegistry' }, ns: 'odh-model-registries' },
    mockModelRegistryList({}),
  );
  cy.interceptOdh(
    'GET /api/service/modelregistry/:serviceName/api/model_registry/:apiVersion/*',
    {
      path: '/registered_models',
      method: 'GET',
    },
    mockRegisteredModelList({}),
  );
  cy.interceptOdh('GET /api/status', mockStatus());
};

describe('Register and Store - Register Model', () => {
  beforeEach(() => {
    initIntercepts();
  });

  it('should show Register and Register and store toggle', () => {
    registerModelPage.visit();
    registerModelPage.findRegisterToggle().should('exist');
    registerModelPage.findRegisterAndStoreToggle().should('exist');
  });

  it('should not show job name field when Register is selected', () => {
    registerModelPage.visit();
    registerModelPage.findRegisterToggle().should('have.attr', 'aria-pressed', 'true');
    cy.get(RegisterModelFormFieldSelector.JOB_NAME).should('not.exist');
  });

  describe('Register and store mode', () => {
    beforeEach(() => {
      registerModelPage.visit();
      registerModelPage.selectRegisterAndStore();
    });

    it('should show job name field when Register and store is selected', () => {
      registerModelPage.findRegisterAndStoreToggle().should('have.attr', 'aria-pressed', 'true');
      cy.get(RegisterModelFormFieldSelector.JOB_NAME).should('exist');
      cy.get(RegisterModelFormFieldSelector.JOB_NAME).should('be.visible');
    });

    it('should allow entering job name', () => {
      const jobName = 'my-transfer-job';
      cy.get(RegisterModelFormFieldSelector.JOB_NAME).type(jobName);
      cy.get(RegisterModelFormFieldSelector.JOB_NAME).should('have.value', jobName);
    });

    it('should show Edit resource name link', () => {
      registerModelPage.findEditResourceNameLink().should('exist');
      registerModelPage.findEditResourceNameLink().should('contain.text', 'Edit resource name');
    });

    it('should show resource name field when Edit resource name is clicked', () => {
      // Resource name field should not be visible initially
      cy.get(RegisterModelFormFieldSelector.RESOURCE_NAME).should('not.exist');

      // Click edit resource name link
      registerModelPage.clickEditResourceName();

      // Resource name field should now be visible
      cy.get(RegisterModelFormFieldSelector.RESOURCE_NAME).should('exist');
      cy.get(RegisterModelFormFieldSelector.RESOURCE_NAME).should('be.visible');
    });

    it('should display helper text for resource name field', () => {
      registerModelPage.clickEditResourceName();

      // Check for validation helper text
      cy.findByText(/Cannot exceed 30 characters/i).should('exist');
      cy.findByText(/Must start and end with a letter or number/i).should('exist');
      cy.findByText(/Auto generated value will be used as resource name if field is blank/i).should(
        'exist',
      );
    });

    it('should allow entering resource name', () => {
      registerModelPage.clickEditResourceName();

      const resourceName = 'my-job-resource';
      cy.get(RegisterModelFormFieldSelector.RESOURCE_NAME).type(resourceName);
      cy.get(RegisterModelFormFieldSelector.RESOURCE_NAME).should('have.value', resourceName);
    });

    it('should show Model origin location section', () => {
      cy.findByText('Model origin location').should('exist');
      cy.findByText('Specify the location that is currently being used to store the model.').should(
        'exist',
      );
    });

    it('should show Model destination location section', () => {
      cy.findByText('Model destination location').should('exist');
      cy.findByText('Specify the location that will be used to store the registered model.').should(
        'exist',
      );
    });

    it('should switch back to Register mode and hide job fields', () => {
      // Verify job field is visible
      cy.get(RegisterModelFormFieldSelector.JOB_NAME).should('exist');

      // Switch back to Register
      registerModelPage.findRegisterToggle().click();

      // Job field should be hidden
      cy.get(RegisterModelFormFieldSelector.JOB_NAME).should('not.exist');
    });
  });
});

describe('Register and Store - Register Version', () => {
  beforeEach(() => {
    initIntercepts();
  });

  it('should show Register and Register and store toggle', () => {
    registerVersionPage.visit('1');
    registerVersionPage.findRegisterToggle().should('exist');
    registerVersionPage.findRegisterAndStoreToggle().should('exist');
  });

  it('should not show job name field when Register is selected', () => {
    registerVersionPage.visit('1');
    registerVersionPage.findRegisterToggle().should('have.attr', 'aria-pressed', 'true');
    cy.get(RegisterVersionFormFieldSelector.JOB_NAME).should('not.exist');
  });

  describe('Register and store mode', () => {
    beforeEach(() => {
      registerVersionPage.visit('1');
      registerVersionPage.selectRegisterAndStore();
    });

    it('should show job name field when Register and store is selected', () => {
      registerVersionPage.findRegisterAndStoreToggle().should('have.attr', 'aria-pressed', 'true');
      cy.get(RegisterVersionFormFieldSelector.JOB_NAME).should('exist');
      cy.get(RegisterVersionFormFieldSelector.JOB_NAME).should('be.visible');
    });

    it('should allow entering job name', () => {
      const jobName = 'version-transfer-job';
      cy.get(RegisterVersionFormFieldSelector.JOB_NAME).type(jobName);
      cy.get(RegisterVersionFormFieldSelector.JOB_NAME).should('have.value', jobName);
    });

    it('should show Edit resource name link', () => {
      registerVersionPage.findEditResourceNameLink().should('exist');
      registerVersionPage.findEditResourceNameLink().should('contain.text', 'Edit resource name');
    });

    it('should show resource name field when Edit resource name is clicked', () => {
      // Resource name field should not be visible initially
      cy.get(RegisterVersionFormFieldSelector.RESOURCE_NAME).should('not.exist');

      // Click edit resource name link
      registerVersionPage.clickEditResourceName();

      // Resource name field should now be visible
      cy.get(RegisterVersionFormFieldSelector.RESOURCE_NAME).should('exist');
      cy.get(RegisterVersionFormFieldSelector.RESOURCE_NAME).should('be.visible');
    });

    it('should display helper text for resource name field', () => {
      registerVersionPage.clickEditResourceName();

      // Check for validation helper text
      cy.findByText(/Cannot exceed 30 characters/i).should('exist');
      cy.findByText(/Must start and end with a letter or number/i).should('exist');
      cy.findByText(/Auto generated value will be used as resource name if field is blank/i).should(
        'exist',
      );
    });

    it('should allow entering resource name', () => {
      registerVersionPage.clickEditResourceName();

      const resourceName = 'version-job-resource';
      cy.get(RegisterVersionFormFieldSelector.RESOURCE_NAME).type(resourceName);
      cy.get(RegisterVersionFormFieldSelector.RESOURCE_NAME).should('have.value', resourceName);
    });
  });
});
