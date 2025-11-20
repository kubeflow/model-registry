import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';
import { BrowserRouter } from 'react-router-dom';
import CreateModal from '~/app/pages/settings/CreateModal';
import * as k8sAPI from '~/app/api/k8s';

// Mock the k8s API
jest.mock('~/app/api/k8s', () => ({
  createModelRegistrySettings: jest.fn(),
}));

// Mock the navigate function
const mockNavigate = jest.fn();
jest.mock('react-router', () => ({
  ...jest.requireActual('react-router'),
  useNavigate: () => mockNavigate,
}));

describe('CreateModal - PostgreSQL Support', () => {
  const mockOnClose = jest.fn();
  const mockRefresh = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    (k8sAPI.createModelRegistrySettings as jest.Mock).mockReturnValue(
      jest.fn().mockResolvedValue({}),
    );
  });

  const renderModal = () =>
    render(
      <BrowserRouter>
        <CreateModal onClose={mockOnClose} refresh={mockRefresh} />
      </BrowserRouter>,
    );

  describe('Database Mode Selection', () => {
    it('should render default database and external database options', () => {
      renderModal();

      expect(screen.getByLabelText('Default database (non-production)')).toBeInTheDocument();
      expect(screen.getByLabelText('External database')).toBeInTheDocument();
    });

    it('should show warning when default database is selected', async () => {
      const user = userEvent.setup();
      renderModal();

      const defaultDatabaseRadio = screen.getByLabelText('Default database (non-production)');
      await user.click(defaultDatabaseRadio);

      expect(
        screen.getByText(/This default database is for development and testing purposes only/i),
      ).toBeInTheDocument();
    });

    it('should show external database fields by default', () => {
      renderModal();

      // External database is the default mode, so fields should be visible immediately
      expect(screen.getByTestId('mr-host-input')).toBeInTheDocument();
      expect(screen.getByTestId('mr-port-input')).toBeInTheDocument();
      expect(screen.getByTestId('mr-username-input')).toBeInTheDocument();
      expect(screen.getByTestId('mr-database-input')).toBeInTheDocument();
    });
  });

  describe('Database Type Selection', () => {
    it('should default to MySQL for external database', () => {
      renderModal();

      // Should show MySQL section title by default
      expect(screen.getByText(/Connect to external MySQL database/i)).toBeInTheDocument();
    });

    it('should allow switching to PostgreSQL', async () => {
      const user = userEvent.setup();
      renderModal();

      // Verify initial state is MySQL
      expect(screen.getByText(/Connect to external MySQL database/i)).toBeInTheDocument();

      // Open database type select
      const databaseTypeToggle = await screen.findByText('MySQL');
      await user.click(databaseTypeToggle);

      // Wait for dropdown to open and find PostgreSQL option
      const postgresOption = await screen.findByText('PostgreSQL');
      await user.click(postgresOption);

      // Verify the section title changed
      await waitFor(() => {
        expect(screen.getByText(/Connect to external PostgreSQL database/i)).toBeInTheDocument();
      });
    });

    it('should auto-fill port when switching database types', async () => {
      const user = userEvent.setup();
      renderModal();

      // Port should default to MySQL port
      const portInput = screen.getByTestId('mr-port-input');
      expect(portInput).toHaveValue('3306');

      // Switch to PostgreSQL
      const databaseTypeToggle = await screen.findByText('MySQL');
      await user.click(databaseTypeToggle);
      const postgresOption = await screen.findByText('PostgreSQL');
      await user.click(postgresOption);

      // Port should auto-fill to PostgreSQL port
      await waitFor(() => {
        expect(portInput).toHaveValue('5432');
      });

      // Switch back to MySQL
      const postgresToggle = await screen.findByText('PostgreSQL');
      await user.click(postgresToggle);
      const mysqlOption = await screen.findByText('MySQL');
      await user.click(mysqlOption);

      // Port should auto-fill back to MySQL port
      await waitFor(() => {
        expect(portInput).toHaveValue('3306');
      });
    });

    it('should not auto-fill port if user manually changed it', async () => {
      const user = userEvent.setup();
      renderModal();

      const portInput = screen.getByTestId('mr-port-input');

      // User manually changes port
      await user.clear(portInput);
      await user.type(portInput, '9999');

      // Switch to PostgreSQL
      const databaseTypeToggle = await screen.findByText('MySQL');
      await user.click(databaseTypeToggle);
      const postgresOption = await screen.findByText('PostgreSQL');
      await user.click(postgresOption);

      // Port should remain user's custom value
      expect(portInput).toHaveValue('9999');
    });
  });

  describe('Form Submission - Default Database', () => {
    it('should create model registry with default PostgreSQL database', async () => {
      const user = userEvent.setup();
      renderModal();

      // Fill in name
      const nameInput = screen.getByTestId('mr-name');
      await user.type(nameInput, 'test-registry');

      // Select default database
      const defaultDatabaseRadio = screen.getByLabelText('Default database (non-production)');
      await user.click(defaultDatabaseRadio);

      // Database field should not be visible for default mode
      expect(screen.queryByTestId('mr-database-input')).not.toBeInTheDocument();

      // Submit
      const createButton = screen.getByTestId('mr-create-button');
      await user.click(createButton);

      await waitFor(() => {
        const mockFn = (k8sAPI.createModelRegistrySettings as jest.Mock).mock.results[0].value;
        expect(mockFn).toHaveBeenCalledWith(
          {},
          expect.objectContaining({
            modelRegistry: expect.objectContaining({
              metadata: expect.objectContaining({
                name: 'test-registry',
              }),
              spec: expect.objectContaining({
                postgres: expect.objectContaining({
                  generateDeployment: true,
                }),
              }),
            }),
          }),
        );
      });
    });
  });

  describe('Form Submission - External MySQL', () => {
    it('should create model registry with external MySQL database', async () => {
      const user = userEvent.setup();
      renderModal();

      // Fill in name
      const nameInput = screen.getByTestId('mr-name');
      await user.type(nameInput, 'test-mysql-registry');

      // Fill in database connection details (external database is selected by default)
      await user.type(screen.getByTestId('mr-host-input'), 'mysql-host');
      // Port is already 3306 by default for MySQL, so no need to change it
      await user.type(screen.getByTestId('mr-username-input'), 'mysql-user');
      const passwordInput = screen.getByTestId('mr-password');
      await user.type(passwordInput, 'mysql-pass');
      await user.type(screen.getByTestId('mr-database-input'), 'model_registry');

      // Submit
      const createButton = screen.getByTestId('mr-create-button');
      await user.click(createButton);

      await waitFor(() => {
        const mockFn = (k8sAPI.createModelRegistrySettings as jest.Mock).mock.results[0].value;
        expect(mockFn).toHaveBeenCalledWith(
          {},
          expect.objectContaining({
            modelRegistry: expect.objectContaining({
              spec: expect.objectContaining({
                mysql: expect.objectContaining({
                  host: 'mysql-host',
                  port: 3306,
                  username: 'mysql-user',
                  database: 'model_registry',
                }),
              }),
            }),
            databasePassword: 'mysql-pass',
          }),
        );
      });
    });
  });

  describe('Form Submission - External PostgreSQL', () => {
    it('should create model registry with external PostgreSQL database', async () => {
      const user = userEvent.setup();
      renderModal();

      // Fill in name
      const nameInput = screen.getByTestId('mr-name');
      await user.type(nameInput, 'test-postgres-registry');

      // External database is selected by default, switch to PostgreSQL
      const databaseTypeToggle = await screen.findByText('MySQL');
      await user.click(databaseTypeToggle);
      const postgresOption = await screen.findByText('PostgreSQL');
      await user.click(postgresOption);

      // Fill in database connection details
      await user.type(screen.getByTestId('mr-host-input'), 'postgres-host');
      // Port should have auto-filled to 5432 when we switched to PostgreSQL
      const portInput = screen.getByTestId('mr-port-input');
      expect(portInput).toHaveValue('5432');
      await user.type(screen.getByTestId('mr-username-input'), 'postgres-user');
      const passwordInput = screen.getByTestId('mr-password');
      await user.type(passwordInput, 'postgres-pass');
      await user.type(screen.getByTestId('mr-database-input'), 'model_registry');

      // Submit
      const createButton = screen.getByTestId('mr-create-button');
      await user.click(createButton);

      await waitFor(() => {
        const mockFn = (k8sAPI.createModelRegistrySettings as jest.Mock).mock.results[0].value;
        expect(mockFn).toHaveBeenCalledWith(
          {},
          expect.objectContaining({
            modelRegistry: expect.objectContaining({
              spec: expect.objectContaining({
                postgres: expect.objectContaining({
                  host: 'postgres-host',
                  port: 5432,
                  username: 'postgres-user',
                  database: 'model_registry',
                }),
              }),
            }),
            databasePassword: 'postgres-pass',
          }),
        );
      });
    });
  });

  describe('Form Validation', () => {
    it('should disable create button when required fields are missing for external database', async () => {
      const user = userEvent.setup();
      renderModal();

      const createButton = screen.getByTestId('mr-create-button');
      expect(createButton).toBeDisabled();

      // Fill only name
      const nameInput = screen.getByTestId('mr-name');
      await user.type(nameInput, 'test-registry');

      // Button should still be disabled
      expect(createButton).toBeDisabled();
    });

    it('should enable create button for default database with minimal fields', async () => {
      const user = userEvent.setup();
      renderModal();

      // Fill in name
      const nameInput = screen.getByTestId('mr-name');
      await user.type(nameInput, 'test-registry');

      // Select default database
      const defaultDatabaseRadio = screen.getByLabelText('Default database (non-production)');
      await user.click(defaultDatabaseRadio);

      // Database field should not be visible for default mode
      expect(screen.queryByTestId('mr-database-input')).not.toBeInTheDocument();

      // Button should be enabled with just name for default mode
      const createButton = screen.getByTestId('mr-create-button');
      await waitFor(() => {
        expect(createButton).not.toBeDisabled();
      });
    });
  });

  describe('Modal Close Behavior', () => {
    it('should call onClose when cancel button is clicked', async () => {
      const user = userEvent.setup();
      renderModal();

      const cancelButton = screen.getByTestId('mr-cancel-button');
      await user.click(cancelButton);

      expect(mockOnClose).toHaveBeenCalled();
    });

    it('should reset form state when modal is closed', async () => {
      const user = userEvent.setup();
      renderModal();

      // Fill in some data
      const nameInput = screen.getByTestId('mr-name');
      await user.type(nameInput, 'test-registry');

      // Close modal
      const cancelButton = screen.getByTestId('mr-cancel-button');
      await user.click(cancelButton);

      expect(mockOnClose).toHaveBeenCalled();
    });
  });
});
