import * as React from 'react';
import {
  PageSection,
  Stack,
  StackItem,
  Button,
  ActionList,
  ActionListItem,
  ActionListGroup,
  Alert,
} from '@patternfly/react-core';
import PreviewButton from './PreviewButton';

type ManageSourceFormFooterProps = {
  submitLabel: string;
  submitError?: Error;
  isSubmitDisabled: boolean;
  isSubmitting: boolean;
  onSubmit: () => void;
  onCancel: () => void;
  isPreviewDisabled: boolean;
  onPreview: () => void;
};

const ManageSourceFormFooter: React.FC<ManageSourceFormFooterProps> = ({
  submitLabel,
  submitError,
  isSubmitDisabled,
  isSubmitting,
  onSubmit,
  onCancel,
  isPreviewDisabled,
  onPreview,
}) => (
  <PageSection hasBodyWrapper={false} stickyOnBreakpoint={{ default: 'bottom' }}>
    <Stack hasGutter>
      {submitError && (
        <StackItem>
          <Alert variant="danger" isInline title="Error saving source">
            {submitError.message}
          </Alert>
        </StackItem>
      )}
      <StackItem>
        <ActionList>
          <ActionListGroup>
            <ActionListItem>
              <Button
                isDisabled={isSubmitDisabled}
                variant="primary"
                id="submit-button"
                data-testid="submit-button"
                isLoading={isSubmitting}
                onClick={onSubmit}
              >
                {submitLabel}
              </Button>
            </ActionListItem>
            <ActionListItem>
              <PreviewButton
                onClick={onPreview}
                isDisabled={isPreviewDisabled}
                variant="secondary"
                testId="preview-button"
              />
            </ActionListItem>
            <ActionListItem>
              <Button
                isDisabled={isSubmitting}
                variant="link"
                id="cancel-button"
                data-testid="cancel-button"
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

export default ManageSourceFormFooter;
