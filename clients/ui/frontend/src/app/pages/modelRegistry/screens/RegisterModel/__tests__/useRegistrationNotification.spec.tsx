import '@testing-library/jest-dom';
import React from 'react';
import { render } from '@testing-library/react';
import { AlertVariant } from '@patternfly/react-core';
import { useRegistrationNotification } from '~/app/pages/modelRegistry/screens/RegisterModel/useRegistrationNotification';

const TITLES = {
  REGISTER_AND_STORE_SUBMITTING: 'Model transfer job started',
  REGISTER_ONLY_SUBMITTING: 'Registering model started',
  REGISTER_ONLY_SUCCESS: 'Model registered',
  REGISTER_ONLY_ERROR: 'Model registration failed',
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
    REGISTER_ONLY_SUBMITTING: 'Registering model started',
    REGISTER_ONLY_SUCCESS: 'Model registered',
    REGISTER_ONLY_ERROR: 'Model registration failed',
  },
}));

jest.mock('~/app/pages/modelRegistry/screens/routeUtils', () => ({
  modelTransferJobsUrl: jest.fn((mrName: string) => `/model-registry/${mrName}/jobs`),
  modelVersionUrl: jest.fn(
    (modelVersionId: string, _registeredModelId: string, mrName: string) =>
      `/model-registry/${mrName}/versions/${modelVersionId}`,
  ),
}));

describe('useRegistrationNotification', () => {
  let setInlineAlert: jest.Mock;

  beforeEach(() => {
    jest.clearAllMocks();
    setInlineAlert = jest.fn();
    mockUseThemeContext.isMUITheme = false;
  });

  function TestWrapper() {
    const actions = useRegistrationNotification(setInlineAlert);
    React.useEffect(() => {
      const toastParams = { versionModelName: 'My Model / v1', mrName: 'mr-sample' };
      actions.showRegisterAndStoreSubmitting(toastParams);
      actions.showRegisterOnlySubmitting();
      actions.showRegisterOnlySuccess({
        ...toastParams,
        modelVersionId: 'mv-1',
        registeredModelId: 'rm-1',
      });
      actions.showRegisterOnlyError({ versionModelName: toastParams.versionModelName });
    }, [actions]);
    return null;
  }

  it('returns all four notification actions', () => {
    function Wrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      expect(actions).toHaveProperty('showRegisterAndStoreSubmitting');
      expect(actions).toHaveProperty('showRegisterOnlySubmitting');
      expect(actions).toHaveProperty('showRegisterOnlySuccess');
      expect(actions).toHaveProperty('showRegisterOnlyError');
      expect(typeof actions.showRegisterAndStoreSubmitting).toBe('function');
      expect(typeof actions.showRegisterOnlySubmitting).toBe('function');
      expect(typeof actions.showRegisterOnlySuccess).toBe('function');
      expect(typeof actions.showRegisterOnlyError).toBe('function');
      return null;
    }
    render(<Wrapper />);
  });

  it('showRegisterAndStoreSubmitting calls notification.info and setInlineAlert when not MUI theme', () => {
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

  it('showRegisterOnlySubmitting calls notification.info with Register only title', () => {
    render(<TestWrapper />);
    expect(mockNotification.info).toHaveBeenCalledWith(
      TITLES.REGISTER_ONLY_SUBMITTING,
      expect.anything(),
    );
    expect(setInlineAlert).toHaveBeenCalledWith(
      expect.objectContaining({
        variant: AlertVariant.info,
        title: TITLES.REGISTER_ONLY_SUBMITTING,
      }),
    );
  });

  it('showRegisterOnlySuccess calls notification.success with Register only success title', () => {
    render(<TestWrapper />);
    expect(mockNotification.success).toHaveBeenCalledWith(
      TITLES.REGISTER_ONLY_SUCCESS,
      expect.anything(),
    );
    expect(setInlineAlert).toHaveBeenCalledWith(
      expect.objectContaining({
        variant: AlertVariant.success,
        title: TITLES.REGISTER_ONLY_SUCCESS,
      }),
    );
  });

  it('showRegisterOnlyError calls notification.error with Register only error title', () => {
    render(<TestWrapper />);
    expect(mockNotification.error).toHaveBeenCalledWith(
      TITLES.REGISTER_ONLY_ERROR,
      expect.anything(),
    );
    expect(setInlineAlert).toHaveBeenCalledWith(
      expect.objectContaining({
        variant: AlertVariant.danger,
        title: TITLES.REGISTER_ONLY_ERROR,
      }),
    );
  });

  it('does not call setInlineAlert when isMUITheme is true', () => {
    mockUseThemeContext.isMUITheme = true;
    function SingleActionWrapper() {
      const actions = useRegistrationNotification(setInlineAlert);
      React.useEffect(() => {
        actions.showRegisterOnlySubmitting();
      }, [actions]);
      return null;
    }
    render(<SingleActionWrapper />);
    expect(mockNotification.info).toHaveBeenCalled();
    expect(setInlineAlert).not.toHaveBeenCalled();
  });
});
