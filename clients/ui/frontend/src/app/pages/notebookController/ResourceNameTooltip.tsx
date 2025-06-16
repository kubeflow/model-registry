import * as React from 'react';
import {
  Box,
  Tooltip,
  Typography,
  IconButton,
} from '@mui/material';
import { HelpOutline } from '@mui/icons-material';
import { K8sResourceCommon } from '@openshift/dynamic-plugin-sdk-utils';
import DashboardPopupIconButton from '~/app/concepts/dashboard/DashboardPopupIconButton';

type ResourceNameTooltipProps = {
  resource: K8sResourceCommon;
  children: React.ReactNode;
};

const ResourceNameTooltip: React.FC<ResourceNameTooltipProps> = ({
  children,
  resource,
}) => (
  <Box sx={{ display: 'inline-flex', alignItems: 'center', gap: 1 }}>
    {children}
    {resource.metadata?.name && (
      <Tooltip
        title={
          <Box>
            <Typography variant="body2">
              Resource names and types are used to find your resources in OpenShift.
            </Typography>
            <Typography variant="body2">
              <b>Resource name:</b> {resource.metadata.name}
            </Typography>
            <Typography variant="body2">
              <b>Resource type:</b> {resource.kind}
            </Typography>
          </Box>
        }
      >
        <DashboardPopupIconButton
          data-testid="resource-name-icon-button"
          aria-label="More info"
        >
            <HelpOutline />
        </DashboardPopupIconButton>
      </Tooltip>
    )}
  </Box>
);

export default ResourceNameTooltip; 