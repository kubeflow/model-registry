import React from 'react';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { ModelRegistryAPIState } from '~/app/hooks/useModelRegistryAPIState';
import useUser from '~/app/hooks/useUser';

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
  const { userId } = useUser();

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
    author: userId,
  };
};
