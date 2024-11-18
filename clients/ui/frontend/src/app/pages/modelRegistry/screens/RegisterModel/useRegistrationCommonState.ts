import React from 'react';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { ModelRegistryAPIState } from '~/app/context/useModelRegistryAPIState';

type RegistrationCommonState = {
  isSubmitting: boolean;
  setIsSubmitting: React.Dispatch<React.SetStateAction<boolean>>;
  submitError: Error | undefined;
  setSubmitError: React.Dispatch<React.SetStateAction<Error | undefined>>;
  handleSubmit: (doSubmit: () => Promise<unknown>) => void;
  apiState: ModelRegistryAPIState;
  author: string;
};

export const useRegistrationCommonState = (): RegistrationCommonState => {
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);

  const { apiState } = React.useContext(ModelRegistryContext);
  const author = 'kubeflow-user'; // TODO: [Auth-enablement] Enable this once we have users ---> useAppSelector((state) => state.user || '');

  const handleSubmit = (doSubmit: () => Promise<unknown>) => {
    setIsSubmitting(true);
    setSubmitError(undefined);
    doSubmit().catch((e: Error) => {
      setIsSubmitting(false);
      setSubmitError(e);
    });
  };

  return {
    isSubmitting,
    setIsSubmitting,
    submitError,
    setSubmitError,
    handleSubmit,
    apiState,
    author,
  };
};
