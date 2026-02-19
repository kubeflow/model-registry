import React, { useRef } from 'react';
import {
  Alert,
  Button,
  FormGroupLabelHelp,
  HelperText,
  HelperTextItem,
  Popover,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { useNamespaceSelector } from 'mod-arch-core';
import ThemeAwareFormGroupWrapper from '~/app/pages/settings/components/ThemeAwareFormGroupWrapper';
import NamespaceSelector from '~/app/standalone/NamespaceSelector';

const NAMESPACE_SELECTOR_TOOLTIP =
  'This list includes only namespaces that you and the selected model registry have permission to access. To request access to a new or existing namespace, contact your administrator.';

const NAMESPACE_NO_ACCESS_MESSAGE =
  'You do not have access to any namespaces. To request access to a new or existing namespace, contact your administrator.';

const SELECTED_NAMESPACE_NO_ACCESS_MESSAGE =
  'The selected namespace does not have access to this model registry. Contact your administrator to grant access.';

const WHO_IS_MY_ADMIN_POPOVER_CONTENT = (
  <Stack hasGutter>
    <StackItem>
      This list includes only namespaces that you have permission to access. To request access to a
      new or existing namespace, contact your administrator.
    </StackItem>
    <StackItem>
      <strong>Your administrator might be:</strong>
    </StackItem>
    <StackItem>
      <ul style={{ margin: 0, paddingLeft: '1.25rem' }}>
        <li>
          The person who assigned you your username, or who helped you log in for the first time
        </li>
        <li>Someone in your IT department or help desk</li>
        <li>A project manager or developer</li>
        <li>Your professor (at a school)</li>
      </ul>
    </StackItem>
  </Stack>
);

export type NamespaceSelectorFieldProps = {
  selectedNamespace: string;
  onSelect: (namespace: string) => void;
  /** Access check result from parent (useCheckNamespaceRegistryAccess); single source of truth to avoid duplicate network calls */
  hasAccess?: boolean | undefined;
  isLoading?: boolean;
  error?: Error | undefined;
};

const NamespaceSelectorField: React.FC<NamespaceSelectorFieldProps> = ({
  selectedNamespace,
  onSelect,
  hasAccess,
  isLoading,
  error,
}) => {
  const labelHelpRef = useRef<HTMLSpanElement>(null);
  const { namespaces = [] } = useNamespaceSelector();
  const isDisabled = namespaces.length === 0;

  const namespaceSelectorElement = (
    <div data-testid="form-namespace-selector-trigger">
      <NamespaceSelector
        isGlobalSelector={false}
        selectedNamespace={selectedNamespace}
        onSelect={onSelect}
        isDisabled={isDisabled}
        placeholderText="Select a namespace"
        isFullWidth
      />
    </div>
  );

  const showNoAccessMessage = namespaces.length === 0;
  const showNoAccessAlert =
    namespaces.length > 0 && selectedNamespace && !isLoading && hasAccess === false;

  const labelHelp = (
    <Popover
      triggerRef={labelHelpRef}
      bodyContent={NAMESPACE_SELECTOR_TOOLTIP}
      aria-label={NAMESPACE_SELECTOR_TOOLTIP}
    >
      <FormGroupLabelHelp ref={labelHelpRef} aria-label="More info for namespace field" />
    </Popover>
  );

  const helperTextNode = (
    <>
      {selectedNamespace && isLoading && (
        <HelperText>
          <HelperTextItem>Checking access...</HelperTextItem>
        </HelperText>
      )}
      {showNoAccessMessage && (
        <Alert
          isInline
          variant="warning"
          title={NAMESPACE_NO_ACCESS_MESSAGE}
          data-testid="namespace-registry-access-alert"
          className="pf-v6-u-mt-sm"
        >
          <Popover bodyContent={WHO_IS_MY_ADMIN_POPOVER_CONTENT} aria-label="Who is my admin?">
            <Button
              variant="link"
              isInline
              component="button"
              data-testid="who-is-my-admin-trigger"
              aria-label="Who is my admin?"
            >
              Who is my admin
            </Button>
          </Popover>
        </Alert>
      )}
      {showNoAccessAlert && (
        <Alert
          isInline
          variant="warning"
          title={SELECTED_NAMESPACE_NO_ACCESS_MESSAGE}
          data-testid="namespace-registry-access-alert"
          className="pf-v6-u-mt-sm"
        />
      )}
      {error && (
        <Alert
          isInline
          variant="danger"
          title="Could not verify namespace access"
          data-testid="namespace-registry-access-error"
          className="pf-v6-u-mt-sm"
        >
          {error.message}
        </Alert>
      )}
    </>
  );

  return (
    <ThemeAwareFormGroupWrapper
      label="Namespace"
      fieldId="namespace-select"
      isRequired
      labelHelp={labelHelp}
      helperTextNode={helperTextNode}
      data-testid="namespace-form-group"
    >
      {namespaceSelectorElement}
    </ThemeAwareFormGroupWrapper>
  );
};

export default NamespaceSelectorField;
