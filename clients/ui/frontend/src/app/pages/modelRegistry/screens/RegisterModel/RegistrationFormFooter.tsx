import React from 'react';
import {
  Alert,
  AlertVariant,
  PageSection,
  Stack,
  StackItem,
  Button,
  ActionList,
  ActionListItem,
  ActionListGroup,
} from '@patternfly/react-core';
import RegisterModelErrors from './RegisterModelErrors';

export type RegistrationInlineAlert = {
  variant: AlertVariant;
  title: string;
  message: React.ReactNode;
};

const FOOTER_STYLE: React.CSSProperties = {
  backgroundColor: 'var(--pf-t--global--background--color--primary--default)',
  zIndex: 100,
  boxShadow: '0 -0.0625rem 0.25rem 0 rgba(3, 3, 3, 0.12)',
};

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
  inlineAlert?: RegistrationInlineAlert;
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
  inlineAlert,
}) => (
  <PageSection
    hasBodyWrapper={false}
    stickyOnBreakpoint={{ default: 'bottom' }}
    style={FOOTER_STYLE}
  >
    <Stack hasGutter>
      {inlineAlert && (
        <StackItem>
          <Alert isInline variant={inlineAlert.variant} title={inlineAlert.title} component="div">
            {inlineAlert.message}
          </Alert>
        </StackItem>
      )}
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
