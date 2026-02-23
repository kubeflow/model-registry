import React from 'react';
import { Link } from 'react-router-dom';
import {
  modelTransferJobsUrl,
  modelVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';

export const REGISTRATION_TOAST_TITLES = {
  REGISTER_AND_STORE_SUBMITTING: 'Model transfer job started',
  REGISTER_AND_STORE_SUCCESS: 'Model transfer job complete',
  REGISTER_AND_STORE_ERROR: 'Model transfer job failed',
  REGISTER_ONLY_SUBMITTING: 'Registering model started',
  REGISTER_ONLY_SUCCESS: 'Model registered',
  REGISTER_ONLY_ERROR: 'Model registration failed',
} as const;

type RegistrationToastMessagesParams = {
  versionModelName: string;
  mrName: string;
  modelVersionId?: string;
  registeredModelId?: string;
};

export const getRegisterAndStoreToastMessageSubmitting = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);

export const getRegisterOnlyToastMessageSubmitting = (): React.ReactNode => 'Please wait.';

export const getRegistrationToastMessageSuccess = ({
  versionModelName,
  mrName,
  modelVersionId,
  registeredModelId,
}: RegistrationToastMessagesParams): React.ReactNode => {
  if (modelVersionId && registeredModelId) {
    return (
      <Link to={modelVersionUrl(modelVersionId, registeredModelId, mrName)}>
        View <strong>{versionModelName}</strong> model version details
      </Link>
    );
  }
  return (
    <>
      View <strong>{versionModelName}</strong> model version details
    </>
  );
};

export const getRegisterAndStoreToastMessageError = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);

export const getRegisterOnlyToastMessageError = ({
  versionModelName,
}: Pick<RegistrationToastMessagesParams, 'versionModelName'>): React.ReactNode => (
  <>
    Registration failed for <strong>{versionModelName}</strong>. Please try again or contact your
    administrator.
  </>
);
