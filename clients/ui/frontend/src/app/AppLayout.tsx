import * as React from 'react';
import { Box } from '@mui/material';
import { Outlet } from 'react-router-dom';

export const AppLayout: React.FC = () => (
  <Box sx={{ display: 'flex' }}>
    <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
      <Outlet />
    </Box>
  </Box>
); 