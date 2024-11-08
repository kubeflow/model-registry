import * as React from 'react';
import {
  sectionTypeBackgroundColor,
  typedBackgroundColor,
  typedObjectImage,
  ProjectObjectType,
  SectionType,
} from '~/shared/components/design/utils';

interface HeaderIconProps {
  size?: number;
  padding?: number;
  image?: string;
  type: ProjectObjectType;
  sectionType?: SectionType;
}

const HeaderIcon: React.FC<HeaderIconProps> = ({
  size = 40,
  padding = 2,
  image,
  type,
  sectionType,
}) => (
  <div
    style={{
      display: 'inline-block',
      width: size,
      height: size,
      padding,
      borderRadius: size / 2,
      background: sectionType
        ? sectionTypeBackgroundColor(sectionType)
        : typedBackgroundColor(type),
    }}
  >
    <img
      width={size - padding * 2}
      height={size - padding * 2}
      src={image || typedObjectImage(type)}
      alt=""
    />
  </div>
);

export default HeaderIcon;
