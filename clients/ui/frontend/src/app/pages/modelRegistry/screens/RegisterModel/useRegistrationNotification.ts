import React from 'react';
import { AlertVariant } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import { useNotification } from '~/app/hooks/useNotification';
import { REGISTRATION_TOAST_TITLES } from '~/app/utilities/const';
import type { RegistrationInlineAlert } from './RegistrationFormFooter';
import {
  getRegisterAndStoreToastMessageSubmitting,
  getRegisterAndStoreToastMessageSuccess,
  getRegisterAndStoreToastMessageError,
} from './registrationToastMessages';

export type RegistrationToastParams = {
  versionModelName: string;
  mrName: string;
};

export type RegistrationNotificationActions = {
  showRegisterAndStoreSubmitting: (params: RegistrationToastParams) => void;
  showRegisterAndStoreSuccess: (params: RegistrationToastParams) => void;
  showRegisterAndStoreError: (params: RegistrationToastParams) => void;
};

/**
 * Shared hook for registration toasts and inline alerts.
 * Shows notification (toast) always; when not using MUI theme, also updates
 * the inline alert in the form footer for consistent UX.
 */
export function useRegistrationNotification(
  setInlineAlert: React.Dispatch<React.SetStateAction<RegistrationInlineAlert | undefined>>,
): RegistrationNotificationActions {
  const notification = useNotification();
  const { isMUITheme } = useThemeContext();

  const showAlert = (variant: AlertVariant, title: string, message: React.ReactNode) => {
    if (variant === AlertVariant.info) {
      notification.info(title, message);
    } else if (variant === AlertVariant.success) {
      notification.success(title, message);
    } else {
      notification.error(title, message);
    }
    if (!isMUITheme) {
      setInlineAlert({ variant, title, message });
    }
  };

  const showRegisterAndStoreSubmitting = (params: RegistrationToastParams) => {
    const title = REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_SUBMITTING;
    const message = getRegisterAndStoreToastMessageSubmitting(params);
    showAlert(AlertVariant.info, title, message);
  };

  const showRegisterAndStoreSuccess = (params: RegistrationToastParams) => {
    const title = REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_SUCCESS;
    const message = getRegisterAndStoreToastMessageSuccess(params);
    showAlert(AlertVariant.success, title, message);
  };

  const showRegisterAndStoreError = (params: RegistrationToastParams) => {
    const title = REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_ERROR;
    const message = getRegisterAndStoreToastMessageError(params);
    showAlert(AlertVariant.danger, title, message);
  };

  return {
    showRegisterAndStoreSubmitting,
    showRegisterAndStoreSuccess,
    showRegisterAndStoreError,
  };
}
