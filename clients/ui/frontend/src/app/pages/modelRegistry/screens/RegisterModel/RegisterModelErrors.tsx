import React from 'react';
import { Alert, AlertActionCloseButton, StackItem } from '@patternfly/react-core';
import { ErrorName, SubmitLabel } from './const';

type RegisterModelErrorProp = {
  submitLabel: string;
  submitError: Error;
  errorName?: string;
  versionName?: string;
  modelName?: string;
};

const RegisterModelErrors: React.FC<RegisterModelErrorProp> = ({
  submitLabel,
  submitError,
  errorName,
  versionName = '',
  modelName = '',
}) => {
  const [showAlert, setShowAlert] = React.useState<boolean>(true);

  if (submitLabel === SubmitLabel.REGISTER_MODEL && errorName === ErrorName.MODEL_VERSION) {
    return (
      <>
        {showAlert && (
          <StackItem>
            <Alert
              isInline
              variant="success"
              title={`${modelName} model registered`}
              actionClose={<AlertActionCloseButton onClose={() => setShowAlert(false)} />}
            />
          </StackItem>
        )}
        <StackItem>
          <Alert isInline variant="danger" title={`Failed to register ${versionName} version`}>
            {submitError.message}
          </Alert>
        </StackItem>
      </>
    );
  }

  if (submitLabel === SubmitLabel.REGISTER_VERSION && errorName === ErrorName.MODEL_VERSION) {
    return (
      <StackItem>
        <Alert isInline variant="danger" title={`Failed to register ${versionName} version`}>
          {submitError.message}
        </Alert>
      </StackItem>
    );
  }

  if (submitLabel === SubmitLabel.REGISTER_MODEL && errorName === ErrorName.MODEL_ARTIFACT) {
    return (
      <>
        {showAlert && (
          <StackItem>
            <Alert
              isInline
              variant="success"
              title={`${modelName} model and ${versionName} version registered`}
              actionClose={<AlertActionCloseButton onClose={() => setShowAlert(false)} />}
            />
          </StackItem>
        )}
        <StackItem>
          <Alert
            isInline
            variant="danger"
            title={`Failed to create artifact for ${versionName} version`}
          >
            {submitError.message}
          </Alert>
        </StackItem>
      </>
    );
  }

  if (submitLabel === SubmitLabel.REGISTER_VERSION && errorName === ErrorName.MODEL_ARTIFACT) {
    return (
      <>
        {showAlert && (
          <StackItem>
            <Alert
              isInline
              variant="success"
              title={`${versionName} version registered`}
              actionClose={<AlertActionCloseButton onClose={() => setShowAlert(false)} />}
            />
          </StackItem>
        )}
        <StackItem>
          <Alert
            isInline
            variant="danger"
            title={`Failed to create artifact for ${versionName} version`}
          >
            {submitError.message}
          </Alert>
        </StackItem>
      </>
    );
  }

  return (
    <StackItem>
      <Alert isInline variant="danger" title={`Failed to register ${modelName} model`}>
        {submitError.message}
      </Alert>
    </StackItem>
  );
};
export default RegisterModelErrors;
