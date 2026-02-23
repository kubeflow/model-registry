import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import {
  REGISTRATION_TOAST_TITLES,
  getRegisterAndStoreToastMessageSubmitting,
  getRegisterAndStoreToastMessageSuccess,
  getRegisterAndStoreToastMessageError,
} from '~/app/pages/modelRegistry/screens/RegisterModel/registrationToastMessages';

jest.mock('~/app/pages/modelRegistry/screens/routeUtils', () => ({
  modelTransferJobsUrl: jest.fn((mrName: string) => `/model-registry/${mrName}/jobs`),
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

  describe('getRegisterAndStoreToastMessageSuccess', () => {
    it('should render message with version name and link to transfer jobs', () => {
      const params = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };
      const node = getRegisterAndStoreToastMessageSuccess(params);
      render(<>{node}</>);

      expect(screen.getByText(/To view/)).toBeInTheDocument();
      expect(screen.getByText('My Model / v1')).toBeInTheDocument();
      const link = screen.getByRole('link', { name: 'Model transfer jobs' });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', '/model-registry/mr-sample/jobs');
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
});
