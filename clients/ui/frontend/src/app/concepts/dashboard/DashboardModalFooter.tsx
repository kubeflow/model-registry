import { Button, Stack, Alert, DialogActions, DialogActionsProps, IconButton } from '@mui/material';
import { Close } from '@mui/icons-material';
import * as React from 'react';

type DashboardModalFooterProps = Pick<DialogActionsProps, 'children'> & {
  onCancel: () => void;
  onSubmit: () => void;
  onReset?: () => void;
  isSubmitLoading?: boolean;
  isSubmitDisabled?: boolean;
  isResetDisabled?: boolean;
  submitLabel?: React.ReactNode;
  resetLabel?: string;
  cancelLabel?: string;
  error?: Error | null;
  alertTitle?: string;
  hideCancel?: boolean;
};

const DashboardModalFooter: React.FC<DashboardModalFooterProps> = ({
  onCancel,
  onSubmit,
  onReset,
  isSubmitLoading,
  isSubmitDisabled,
  isResetDisabled,
  submitLabel,
  resetLabel = 'Reset',
  cancelLabel = 'Cancel',
  error,
  alertTitle,
  hideCancel,
  children,
}) => (
  <DialogActions>
    <Stack spacing={2}>
      {error && (
        <Alert
          severity="error"
          action={
            <IconButton size="small" aria-label="close" color="inherit" onClick={onReset}>
              <Close fontSize="small" />
            </IconButton>
          }
        >
          {alertTitle && <b>{alertTitle}</b>}
          {error.message}
        </Alert>
      )}
      <Stack direction="row" spacing={1}>
        <Button
          variant="contained"
          onClick={onSubmit}
          disabled={isSubmitDisabled || isSubmitLoading}
        >
          {isSubmitLoading ? 'Loading...' : submitLabel}
        </Button>
        {onReset && (
          <Button variant="text" onClick={onReset} disabled={isResetDisabled}>
            {resetLabel}
          </Button>
        )}
        {!hideCancel && (
          <Button variant="outlined" onClick={onCancel}>
            {cancelLabel}
          </Button>
        )}
      </Stack>
      {children}
    </Stack>
  </DialogActions>
);

export default DashboardModalFooter;
