import '@testing-library/jest-dom';
import React from 'react';
import { render, screen } from '@testing-library/react';
import { getRegisterAndStoreToastMessage } from '~/app/pages/modelRegistry/screens/RegisterModel/registrationToastMessages';

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
  describe('getRegisterAndStoreToastMessage', () => {
    it('should render message with version name and link to transfer jobs', () => {
      const params = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };
      const node = getRegisterAndStoreToastMessage(params);
      render(<>{node}</>);

      expect(screen.getByText(/To view/)).toBeInTheDocument();
      expect(screen.getByText('My Model / v1')).toBeInTheDocument();
      const link = screen.getByRole('link', { name: 'Model transfer jobs' });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', '/model-registry/mr-sample/jobs');
    });
  });
});
