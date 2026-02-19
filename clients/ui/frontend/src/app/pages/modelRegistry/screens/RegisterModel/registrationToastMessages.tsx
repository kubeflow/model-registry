import React from 'react';
import { Link } from 'react-router-dom';
import {
  modelTransferJobsUrl,
  modelVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';

export const REGISTRATION_TOAST_TITLES = {
  SUBMITTING: 'Model transfer job started',
  SUCCESS: 'Model transfer job complete',
  ERROR: 'Model transfer job failed',
} as const;

type RegistrationToastMessagesParams = {
  versionModelName: string;
  mrName: string;
  modelVersionId?: string;
  registeredModelId?: string;
};

export const getRegistrationToastMessageSubmitting = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);

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

export const getRegistrationToastMessageError = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);
