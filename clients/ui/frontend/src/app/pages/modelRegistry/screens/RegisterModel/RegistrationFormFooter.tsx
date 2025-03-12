import React from 'react';
import {
  PageSection,
  Stack,
  StackItem,
  Button,
  ActionList,
  ActionListItem,
  ActionListGroup,
} from '@patternfly/react-core';
import RegisterModelErrors from './RegisterModelErrors';

type RegistrationFormFooterProps = {
  submitLabel: string;
  submitError?: Error;
  isSubmitDisabled: boolean;
  isSubmitting: boolean;
  onSubmit: () => void;
  onCancel: () => void;
  registrationErrorType?: string;
  versionName?: string;
  modelName?: string;
};

const RegistrationFormFooter: React.FC<RegistrationFormFooterProps> = ({
  submitLabel,
  submitError,
  isSubmitDisabled,
  isSubmitting,
  onSubmit,
  onCancel,
  registrationErrorType,
  versionName,
  modelName,
}) => (
  <PageSection hasBodyWrapper={false} stickyOnBreakpoint={{ default: 'bottom' }}>
    <Stack hasGutter>
      {submitError && (
        <RegisterModelErrors
          submitLabel={submitLabel}
          submitError={submitError}
          registrationErrorType={registrationErrorType}
          versionName={versionName}
          modelName={modelName}
        />
      )}
      <StackItem>
        <ActionList>
          <ActionListGroup>
            <ActionListItem>
              <Button
                isDisabled={isSubmitDisabled}
                variant="primary"
                id="create-button"
                data-testid="create-button"
                isLoading={isSubmitting}
                onClick={onSubmit}
              >
                {submitLabel}
              </Button>
            </ActionListItem>
            <ActionListItem>
              <Button
                isDisabled={isSubmitting}
                variant="link"
                id="cancel-button"
                onClick={onCancel}
              >
                Cancel
              </Button>
            </ActionListItem>
          </ActionListGroup>
        </ActionList>
      </StackItem>
    </Stack>
  </PageSection>
);

export default RegistrationFormFooter;
