import React from 'react';
import { AlertVariant } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import { useNotification, type NotificationLinkOptions } from '~/app/hooks/useNotification';
import { REGISTRATION_TOAST_TITLES } from '~/app/utilities/const';
import { modelTransferJobsUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import type { RegistrationInlineAlert } from './RegistrationFormFooter';
import {
  getRegisterAndStoreToastMessage,
  type RegistrationToastMessagesParams,
} from './registrationToastMessages';

export type RegistrationNotificationActions = {
  showRegisterAndStoreStarted: (params: RegistrationToastMessagesParams) => void;
  showRegisterAndStoreError: (params: RegistrationToastMessagesParams) => void;
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

  const showAlert = (
    variant: AlertVariant,
    title: string,
    message: React.ReactNode,
    options?: NotificationLinkOptions,
  ) => {
    if (variant === AlertVariant.info) {
      notification.info(title, message, options);
    } else if (variant === AlertVariant.success) {
      notification.success(title, message, options);
    } else {
      notification.error(title, message, options);
    }
    if (!isMUITheme) {
      setInlineAlert({ variant, title, message });
    }
  };

  const showRegisterAndStoreStarted = (params: RegistrationToastMessagesParams) => {
    const title = REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_STARTED;
    const message = getRegisterAndStoreToastMessage(params);
    showAlert(AlertVariant.info, title, message, {
      linkUrl: modelTransferJobsUrl(params.mrName),
      linkLabel: 'Model transfer jobs',
      messageText: `To view ${params.versionModelName} job details, go to `,
    });
  };

  const showRegisterAndStoreError = (params: RegistrationToastMessagesParams) => {
    const title = REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_ERROR;
    const message = getRegisterAndStoreToastMessage(params);
    showAlert(AlertVariant.danger, title, message);
  };

  return {
    showRegisterAndStoreStarted,
    showRegisterAndStoreError,
  };
}
