import '@testing-library/jest-dom';
import React from 'react';
import { render } from '@testing-library/react';
import { AlertVariant } from '@patternfly/react-core';
import { useRegistrationNotification } from '~/app/pages/modelRegistry/screens/RegisterModel/useRegistrationNotification';

const TITLES = {
  REGISTER_AND_STORE_SUBMITTING: 'Model transfer job started',
  REGISTER_AND_STORE_SUCCESS: 'Model transfer job complete',
  REGISTER_AND_STORE_ERROR: 'Model transfer job failed',
};

const mockNotification = {
  info: jest.fn(),
  success: jest.fn(),
  error: jest.fn(),
};

jest.mock('~/app/hooks/useNotification', () => ({
  useNotification: () => mockNotification,
}));

const mockUseThemeContext = { isMUITheme: false };
jest.mock('mod-arch-kubeflow', () => ({
  useThemeContext: () => mockUseThemeContext,
}));

jest.mock('~/app/utilities/const', () => ({
  REGISTRATION_TOAST_TITLES: {
    REGISTER_AND_STORE_SUBMITTING: 'Model transfer job started',
    REGISTER_AND_STORE_SUCCESS: 'Model transfer job complete',
    REGISTER_AND_STORE_ERROR: 'Model transfer job failed',
  },
}));

jest.mock('~/app/pages/modelRegistry/screens/routeUtils', () => ({
  modelTransferJobsUrl: jest.fn((mrName: string) => `/model-registry/${mrName}/jobs`),
}));

describe('useRegistrationNotification', () => {
  let setInlineAlert: jest.Mock;
  const toastParams = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };

  beforeEach(() => {
    jest.clearAllMocks();
    setInlineAlert = jest.fn();
    mockUseThemeContext.isMUITheme = false;
  });

  it('returns Register and Store notification actions only', () => {
    function Wrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      expect(actions).toHaveProperty('showRegisterAndStoreSubmitting');
      expect(actions).toHaveProperty('showRegisterAndStoreSuccess');
      expect(actions).toHaveProperty('showRegisterAndStoreError');
      expect(typeof actions.showRegisterAndStoreSubmitting).toBe('function');
      expect(typeof actions.showRegisterAndStoreSuccess).toBe('function');
      expect(typeof actions.showRegisterAndStoreError).toBe('function');
      return null;
    }
    render(<Wrapper />);
  });

  it('showRegisterAndStoreSubmitting calls notification.info and setInlineAlert when not MUI theme', () => {
    function TestWrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      React.useEffect(() => {
        actions.showRegisterAndStoreSubmitting(toastParams);
      }, [actions]);
      return null;
    }
    render(<TestWrapper />);
    expect(mockNotification.info).toHaveBeenCalledWith(
      TITLES.REGISTER_AND_STORE_SUBMITTING,
      expect.anything(),
    );
    expect(setInlineAlert).toHaveBeenCalledWith(
      expect.objectContaining({
        variant: AlertVariant.info,
        title: TITLES.REGISTER_AND_STORE_SUBMITTING,
      }),
    );
  });

  it('showRegisterAndStoreSuccess calls notification.success and setInlineAlert when not MUI theme', () => {
    function TestWrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      React.useEffect(() => {
        actions.showRegisterAndStoreSuccess(toastParams);
      }, [actions]);
      return null;
    }
    render(<TestWrapper />);
    expect(mockNotification.success).toHaveBeenCalledWith(
      TITLES.REGISTER_AND_STORE_SUCCESS,
      expect.anything(),
    );
    expect(setInlineAlert).toHaveBeenCalledWith(
      expect.objectContaining({
        variant: AlertVariant.success,
        title: TITLES.REGISTER_AND_STORE_SUCCESS,
      }),
    );
  });

  it('showRegisterAndStoreError calls notification.error and setInlineAlert when not MUI theme', () => {
    function TestWrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      React.useEffect(() => {
        actions.showRegisterAndStoreError(toastParams);
      }, [actions]);
      return null;
    }
    render(<TestWrapper />);
    expect(mockNotification.error).toHaveBeenCalledWith(
      TITLES.REGISTER_AND_STORE_ERROR,
      expect.anything(),
    );
    expect(setInlineAlert).toHaveBeenCalledWith(
      expect.objectContaining({
        variant: AlertVariant.danger,
        title: TITLES.REGISTER_AND_STORE_ERROR,
      }),
    );
  });

  it('does not call setInlineAlert when isMUITheme is true', () => {
    mockUseThemeContext.isMUITheme = true;
    function TestWrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      React.useEffect(() => {
        actions.showRegisterAndStoreSubmitting(toastParams);
      }, [actions]);
      return null;
    }
    render(<TestWrapper />);
    expect(mockNotification.info).toHaveBeenCalled();
    expect(setInlineAlert).not.toHaveBeenCalled();
  });
});
