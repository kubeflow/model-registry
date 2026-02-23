import React from 'react';
import { Link } from 'react-router-dom';
import { REGISTRATION_TOAST_TITLES } from '~/app/utilities/const';
import { modelTransferJobsUrl } from '~/app/pages/modelRegistry/screens/routeUtils';

export { REGISTRATION_TOAST_TITLES };

type RegistrationToastMessagesParams = {
  versionModelName: string;
  mrName: string;
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

export const getRegisterAndStoreToastMessageSuccess = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);

export const getRegisterAndStoreToastMessageError = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);
