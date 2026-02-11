import React, { useRef } from 'react';
import { Link } from 'react-router-dom';
import {
  Box,
  FormControl,
  IconButton,
  MenuItem,
  Select,
  SelectChangeEvent,
  Tooltip,
} from '@mui/material';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import { Alert, FormGroup, FormGroupLabelHelp, Popover } from '@patternfly/react-core';
import { useNamespaceSelector } from 'mod-arch-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import { useCheckNamespaceRegistryAccess } from '~/app/hooks/useCheckNamespaceRegistryAccess';

const NAMESPACE_SELECTOR_TOOLTIP =
  'This list includes only projects that you and the selected model registry have permission to access. To request access to a new or existing project, contact your administrator.';

const NAMESPACE_NO_ACCESS_MESSAGE =
  'The selected namespace does not have access to this model registry. Please contact your administrator to grant access or select a different namespace.';

export type NamespaceSelectorFieldProps = {
  selectedNamespace: string;
  onSelect: (namespace: string) => void;
  registryName?: string;
  registryNamespace?: string;
  onAccessChange?: (hasAccess: boolean | undefined) => void;
};

/**
 * Namespace selector field with registry access validation (SSAR).
 * Renders a FormGroup with label help, MUI Select for namespace choice, and optional alerts for no access / error.
 * Use in register-and-store flows; syncs access result to parent via onAccessChange.
 */
const NamespaceSelectorField: React.FC<NamespaceSelectorFieldProps> = ({
  selectedNamespace,
  onSelect,
  registryName,
  registryNamespace,
  onAccessChange,
}) => {
  const labelHelpRef = useRef<HTMLSpanElement>(null);
  const { hasAccess, isLoading, error } = useCheckNamespaceRegistryAccess(
    registryName,
    registryNamespace,
    selectedNamespace,
  );

  React.useEffect(() => {
    onAccessChange?.(hasAccess);
  }, [hasAccess, onAccessChange]);

  const { namespaces = [] } = useNamespaceSelector();
  const { isMUITheme } = useThemeContext();
  const isDisabled = namespaces.length === 0;

  const handleChange = (event: SelectChangeEvent<string>) => {
    const { value } = event.target;
    if (value) {
      onSelect(value);
    }
  };

  const selectControl = (
    <FormControl fullWidth size="medium" disabled={isDisabled} required>
      <Select
        displayEmpty
        value={selectedNamespace || ''}
        onChange={handleChange}
        renderValue={(value) => value || 'Select a namespace'}
      >
        <MenuItem value="" disabled>
          Select a namespace
        </MenuItem>
        {namespaces.map((ns) => (
          <MenuItem key={ns.name} value={ns.name}>
            {ns.name}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );

  const namespaceSelectorElement = (
    <div data-testid="form-namespace-selector">
      {isMUITheme ? (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <div data-testid="form-namespace-selector-trigger" style={{ flex: 1, minWidth: 0 }}>
            {selectControl}
          </div>
          <Tooltip title={NAMESPACE_SELECTOR_TOOLTIP} placement="top">
            <IconButton
              size="small"
              aria-label="More info for namespace field"
              sx={{ flexShrink: 0 }}
            >
              <HelpOutlineIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      ) : (
        <div data-testid="form-namespace-selector-trigger">{selectControl}</div>
      )}
    </div>
  );

  const showNoAccessAlert = selectedNamespace && !isLoading && hasAccess === false;

  return (
    <FormGroup
      label="Namespace"
      fieldId="namespace-select"
      isRequired
      data-testid="namespace-form-group"
      labelHelp={
        !isMUITheme ? (
          <Popover
            triggerRef={labelHelpRef}
            bodyContent={NAMESPACE_SELECTOR_TOOLTIP}
            aria-label={NAMESPACE_SELECTOR_TOOLTIP}
          >
            <FormGroupLabelHelp ref={labelHelpRef} aria-label="More info for namespace field" />
          </Popover>
        ) : undefined
      }
    >
      {namespaceSelectorElement}
      {showNoAccessAlert && (
        <Alert
          isInline
          variant="warning"
          title={NAMESPACE_NO_ACCESS_MESSAGE}
          data-testid="namespace-registry-access-alert"
          className="pf-v6-u-mt-sm"
        >
          {registryName && (
            <Link to="/model-registry-settings" target="_blank" rel="noopener noreferrer">
              Configure namespace access in Model Registry settings
            </Link>
          )}
        </Alert>
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
    </FormGroup>
  );
};

export default NamespaceSelectorField;
