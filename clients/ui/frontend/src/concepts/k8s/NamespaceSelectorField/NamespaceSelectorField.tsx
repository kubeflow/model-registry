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
import { NamespaceSelectorMessages } from '~/app/utilities/const';

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
      <ul>
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
      bodyContent={NamespaceSelectorMessages.SELECTOR_TOOLTIP}
      aria-label={NamespaceSelectorMessages.SELECTOR_TOOLTIP}
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
          title={NamespaceSelectorMessages.NO_ACCESS}
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
          title={NamespaceSelectorMessages.SELECTED_NAMESPACE_NO_ACCESS}
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
