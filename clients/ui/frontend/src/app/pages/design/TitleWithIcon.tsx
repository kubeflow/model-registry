import * as React from 'react';
import { Box, Typography } from '@mui/material';
import {
  ProjectObjectType,
  typedBackgroundColor,
  typedObjectImage,
} from '~/app/concepts/design/utils';
import TypedObjectIcon from '~/app/concepts/design/TypedObjectIcon';

interface TitleWithIconProps {
  title: React.ReactNode;
  objectType: ProjectObjectType;
  iconSize?: number;
  padding?: number;
}

const TitleWithIcon: React.FC<TitleWithIconProps> = ({
  title,
  objectType,
  iconSize = 40,
  padding = 4,
}) => (
  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
    <Box
      sx={{
        background: typedBackgroundColor(objectType),
        borderRadius: '50%',
        padding: `${padding}px`,
        width: iconSize,
        height: iconSize,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <TypedObjectIcon
        resourceType={objectType}
        style={{ width: iconSize - padding * 2, height: iconSize - padding * 2 }}
        src={typedObjectImage(objectType)}
      />
    </Box>
    <Typography variant="h6">{title}</Typography>
  </Box>
);

export default TitleWithIcon; 