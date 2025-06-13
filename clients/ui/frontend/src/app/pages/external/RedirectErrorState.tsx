import React from 'react';
import { Box, Typography } from '@mui/material';
import { Error } from '@mui/icons-material';

type RedirectErrorStateProps = {
  title?: string;
  errorMessage?: string;
  actions?: React.ReactNode | React.ReactNode[];
};

const RedirectErrorState: React.FC<RedirectErrorStateProps> = ({
  title,
  errorMessage,
  actions,
}) => (
  <Box sx={{ textAlign: 'center', py: 4 }}>
    <Error color="error" sx={{ fontSize: 48 }} />
    <Typography variant="h5" component="h1" sx={{ mt: 2 }}>
      {title ?? 'Error redirecting'}
    </Typography>
    {errorMessage && (
      <Typography variant="body1" sx={{ mt: 1 }}>
        {errorMessage}
      </Typography>
    )}
    {actions && <Box sx={{ mt: 2 }}>{actions}</Box>}
  </Box>
);

export default RedirectErrorState; 