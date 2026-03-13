import React from 'react';
import { Link } from 'react-router-dom';
import { modelTransferJobsUrl } from '~/app/pages/modelRegistry/screens/routeUtils';

export type RegistrationToastMessagesParams = {
  versionModelName: string;
  mrName: string;
};

export const getRegisterAndStoreToastMessage = ({
  versionModelName,
  mrName,
}: RegistrationToastMessagesParams): React.ReactNode => (
  <>
    To view <strong>{versionModelName}</strong> job details, go to{' '}
    <Link to={modelTransferJobsUrl(mrName)}>Model transfer jobs</Link>.
  </>
);
