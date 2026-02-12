import React, { useRef } from 'react';
import { Link } from 'react-router-dom';
import {
  Box,
  FormControl,
  IconButton,
  MenuItem,
  Select,
  SelectChangeEvent,
  Typography,
} from '@mui/material';
import ErrorIcon from '@mui/icons-material/Error';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import {
  Alert,
  FormGroup,
  FormGroupLabelHelp,
  Popover,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { useNamespaceSelector } from 'mod-arch-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import { useCheckNamespaceRegistryAccess } from '~/app/hooks/useCheckNamespaceRegistryAccess';
import useUser from '~/app/hooks/useUser';

const NAMESPACE_SELECTOR_TOOLTIP =
  'This list includes only namespaces that you and the selected model registry have permission to access. To request access to a new or existing namespace, contact your administrator.';

const NAMESPACE_NO_ACCESS_MESSAGE =
  'You do not have access to any namespaces. To request access to a new or existing namespace, contact your administrator.';

const SELECTED_NAMESPACE_NO_ACCESS_MESSAGE_ADMIN =
  'The selected namespace does not have access to this model registry. To grant access, click Manage permissions for this registry on the Model registry settings page.';
const SELECTED_NAMESPACE_NO_ACCESS_MESSAGE_USER =
  'The selected namespace does not have access to this model registry. Contact your administrator to grant access.';
const MODEL_REGISTRY_SETTINGS_LINK_LABEL = 'Go to Model registry settings';

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
  registryName?: string;
  registryNamespace?: string;
  onAccessChange?: (hasAccess: boolean | undefined) => void;
};


const NamespaceSelectorField: React.FC<NamespaceSelectorFieldProps> = ({
  selectedNamespace,
  onSelect,
  registryName,
  registryNamespace,
  onAccessChange,
}) => {
  const labelHelpRef = useRef<HTMLSpanElement>(null);
  const { clusterAdmin } = useUser();
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
          <Popover bodyContent={NAMESPACE_SELECTOR_TOOLTIP} aria-label={NAMESPACE_SELECTOR_TOOLTIP}>
            <IconButton
              size="small"
              aria-label="More info for namespace field"
              sx={{ flexShrink: 0 }}
            >
              <HelpOutlineIcon fontSize="small" />
            </IconButton>
          </Popover>
        </Box>
      ) : (
        <div data-testid="form-namespace-selector-trigger">{selectControl}</div>
      )}
    </div>
  );

  const showNoAccessMessage = namespaces.length === 0;
  const showNoAccessAlert =
    namespaces.length > 0 && selectedNamespace && !isLoading && hasAccess === false;

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
      {showNoAccessMessage && (
        <Box
          data-testid="namespace-registry-access-alert"
          sx={{
            display: 'flex',
            alignItems: 'flex-start',
            gap: 1,
            mt: 1.5,
          }}
        >
          <ErrorIcon sx={{ color: 'error.main', fontSize: 20, flexShrink: 0 }} aria-hidden />
          <Typography variant="body2" component="span" sx={{ flex: 1, minWidth: 0 }}>
            {NAMESPACE_NO_ACCESS_MESSAGE}{' '}
            <Popover bodyContent={WHO_IS_MY_ADMIN_POPOVER_CONTENT} aria-label="Who is my admin?">
              <Typography
                component="button"
                type="button"
                variant="body2"
                data-testid="who-is-my-admin-trigger"
                sx={{
                  display: 'inline',
                  p: 0,
                  border: 'none',
                  background: 'none',
                  cursor: 'pointer',
                  color: 'primary.main',
                  textDecoration: 'underline',
                  font: 'inherit',
                  '&:hover': { color: 'primary.dark' },
                }}
                aria-label="Who is my admin?"
              >
                <HelpOutlineIcon
                  sx={{ fontSize: 16, verticalAlign: 'middle', mr: 0.25 }}
                  aria-hidden
                />
                Who is my admin
              </Typography>
            </Popover>
          </Typography>
        </Box>
      )}
      {showNoAccessAlert && (
        <Alert
          isInline
          variant="warning"
          title={
            clusterAdmin
              ? SELECTED_NAMESPACE_NO_ACCESS_MESSAGE_ADMIN
              : SELECTED_NAMESPACE_NO_ACCESS_MESSAGE_USER
          }
          data-testid="namespace-registry-access-alert"
          className="pf-v6-u-mt-sm"
        >
          {clusterAdmin && (
            <Link to="/model-registry-settings" target="_blank" rel="noopener noreferrer">
              {MODEL_REGISTRY_SETTINGS_LINK_LABEL}
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
