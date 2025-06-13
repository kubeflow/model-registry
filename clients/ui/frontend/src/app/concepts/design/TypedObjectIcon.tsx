import * as React from 'react';
import { SvgIcon } from '@mui/material';
import { ProjectObjectType } from '~/app/concepts/design/utils';

type TypedObjectIconProps = {
  resourceType: ProjectObjectType;
  style?: React.CSSProperties;
  src?: string;
};

const TypedObjectIcon: React.FC<TypedObjectIconProps> = ({ resourceType, style, src }) => (
  <SvgIcon style={style}>
    <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z" />
  </SvgIcon>
);

export default TypedObjectIcon; 