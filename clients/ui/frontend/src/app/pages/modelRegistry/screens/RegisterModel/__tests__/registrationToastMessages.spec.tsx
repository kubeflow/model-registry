import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import {
  REGISTRATION_TOAST_TITLES,
  getRegisterAndStoreToastMessageSubmitting,
  getRegisterOnlyToastMessageSubmitting,
  getRegistrationToastMessageSuccess,
  getRegisterAndStoreToastMessageError,
  getRegisterOnlyToastMessageError,
} from '~/app/pages/modelRegistry/screens/RegisterModel/registrationToastMessages';

jest.mock('~/app/pages/modelRegistry/screens/routeUtils', () => ({
  modelTransferJobsUrl: jest.fn((mrName: string) => `/model-registry/${mrName}/jobs`),
  modelVersionUrl: jest.fn(
    (modelVersionId: string, _registeredModelId: string, mrName: string) =>
      `/model-registry/${mrName}/versions/${modelVersionId}`,
  ),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  Link: ({ to, children, ...props }: { to: string; children: React.ReactNode }) => (
    <a href={to} {...props}>
      {children}
    </a>
  ),
}));

describe('registrationToastMessages', () => {
  describe('REGISTRATION_TOAST_TITLES', () => {
    it('should have expected title constants for Register and Store flow', () => {
      expect(REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_SUBMITTING).toBe(
        'Model transfer job started',
      );
      expect(REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_SUCCESS).toBe(
        'Model transfer job complete',
      );
      expect(REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_ERROR).toBe('Model transfer job failed');
    });

    it('should have expected title constants for Register only flow', () => {
      expect(REGISTRATION_TOAST_TITLES.REGISTER_ONLY_SUBMITTING).toBe(
        'Registering version started',
      );
      expect(REGISTRATION_TOAST_TITLES.REGISTER_ONLY_SUCCESS).toBe('Version registered');
      expect(REGISTRATION_TOAST_TITLES.REGISTER_ONLY_ERROR).toBe('Version registration failed');
    });
  });

  describe('getRegisterAndStoreToastMessageSubmitting', () => {
    it('should render message with version name and link to transfer jobs', () => {
      const params = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };
      const node = getRegisterAndStoreToastMessageSubmitting(params);
      render(<>{node}</>);

      expect(screen.getByText(/To view/)).toBeInTheDocument();
      expect(screen.getByText('My Model / v1')).toBeInTheDocument();
      const link = screen.getByRole('link', { name: 'Model transfer jobs' });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', '/model-registry/mr-sample/jobs');
    });
  });

  describe('getRegisterOnlyToastMessageSubmitting', () => {
    it('should return please wait message', () => {
      const node = getRegisterOnlyToastMessageSubmitting();
      render(<>{node}</>);
      expect(screen.getByText('Please wait.')).toBeInTheDocument();
    });
  });

  describe('getRegistrationToastMessageSuccess', () => {
    it('should return Link to version details when modelVersionId and registeredModelId are provided', () => {
      const params = {
        versionModelName: 'My Model / v1',
        mrName: 'mr-sample',
        modelVersionId: 'mv-123',
        registeredModelId: 'rm-456',
      };
      const node = getRegistrationToastMessageSuccess(params);
      render(<>{node}</>);

      const link = screen.getByRole('link', {
        name: /View.*My Model \/ v1.*model version details/,
      });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', '/model-registry/mr-sample/versions/mv-123');
    });

    it('should return fragment with text when modelVersionId or registeredModelId is missing', () => {
      const params = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };
      const node = getRegistrationToastMessageSuccess(params);
      render(<>{node}</>);

      expect(screen.getByText(/View/)).toBeInTheDocument();
      expect(screen.getByText('My Model / v1')).toBeInTheDocument();
      expect(screen.getByText(/model version details/)).toBeInTheDocument();
      expect(screen.queryByRole('link')).not.toBeInTheDocument();
    });
  });

  describe('getRegisterAndStoreToastMessageError', () => {
    it('should render message with version name and link to transfer jobs', () => {
      const params = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };
      const node = getRegisterAndStoreToastMessageError(params);
      render(<>{node}</>);

      expect(screen.getByText(/To view/)).toBeInTheDocument();
      expect(screen.getByText('My Model / v1')).toBeInTheDocument();
      const link = screen.getByRole('link', { name: 'Model transfer jobs' });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', '/model-registry/mr-sample/jobs');
    });
  });

  describe('getRegisterOnlyToastMessageError', () => {
    it('should render message with version name and no transfer job link', () => {
      const params = { versionModelName: 'My Model / v1' };
      const node = getRegisterOnlyToastMessageError(params);
      render(<>{node}</>);

      expect(screen.getByText(/Registration failed for/)).toBeInTheDocument();
      expect(screen.getByText('My Model / v1')).toBeInTheDocument();
      expect(screen.getByText(/Please try again or contact your administrator/)).toBeInTheDocument();
      expect(screen.queryByRole('link')).not.toBeInTheDocument();
    });
  });
});
