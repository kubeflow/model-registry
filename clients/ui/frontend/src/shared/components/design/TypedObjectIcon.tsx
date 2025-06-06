import * as React from 'react';
import { SVGIconProps } from '@patternfly/react-icons/dist/esm/createIcon';
import { ProjectObjectType, typedColor } from '~/shared/components/design/utils';
import { GroupIcon, RegisteredModelIcon, UserIcon } from '~/shared/images/icons';

type TypedObjectIconProps = SVGIconProps & {
  resourceType: ProjectObjectType;
  useTypedColor?: boolean;
  size?: number;
};
const TypedObjectIcon: React.FC<TypedObjectIconProps> = ({
  resourceType,
  useTypedColor,
  style,
  ...rest
}) => {
  let Icon;

  switch (resourceType) {
    case ProjectObjectType.registeredModels:
    case ProjectObjectType.modelRegistrySettings:
      Icon = RegisteredModelIcon;
      break;
    case ProjectObjectType.user:
      Icon = UserIcon;
      break;
    case ProjectObjectType.group:
      Icon = GroupIcon;
      break;
    default:
      return null;
  }

  return (
    <Icon
      style={{
        color: useTypedColor ? typedColor(resourceType) : undefined,
        ...(style || {}),
      }}
      {...rest}
    />
  );
};

export default TypedObjectIcon;
