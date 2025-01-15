import * as React from 'react';
import { Flex, FlexItem } from '@patternfly/react-core';
import {
  ProjectObjectType,
  typedBackgroundColor,
  typedObjectImage,
} from '~/shared/components/design/utils';
import TypedObjectIcon from '~/shared/components/design/TypedObjectIcon';

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
  <Flex spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsCenter' }}>
    <FlexItem>
      <div
        style={{
          background: typedBackgroundColor(objectType),
          borderRadius: iconSize / 2,
          padding,
          width: iconSize,
          height: iconSize,
        }}
      >
        <TypedObjectIcon
          resourceType={objectType}
          style={{ width: iconSize - padding * 2, height: iconSize - padding * 2 }}
          src={typedObjectImage(objectType)}
        />
      </div>
    </FlexItem>
    <FlexItem>{title}</FlexItem>
  </Flex>
);

export default TitleWithIcon;
