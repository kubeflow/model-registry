import * as React from 'react';
import { Box, Typography, Divider } from '@mui/material';

type FormSectionProps = {
  title: string;
  description?: string;
  children: React.ReactNode;
};

const FormSection: React.FC<FormSectionProps> = ({ title, description, children }) => (
  <Box sx={{ my: 2 }}>
    <Typography variant="h6">{title}</Typography>
    {description && <Typography variant="body2" color="text.secondary">{description}</Typography>}
    <Box sx={{ mt: 2 }}>{children}</Box>
    <Divider sx={{ mt: 2 }} />
  </Box>
);

export default FormSection; 