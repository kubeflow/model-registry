import * as React from 'react';
import { DescriptionListGroup, DescriptionListDescription } from '@patternfly/react-core';
import { ModelRegistryCustomProperties } from '~/app/types';
import ModelPropertiesExpandableSection from '~/app/pages/modelRegistry/screens/components/ModelPropertiesExpandableSection';

type ModelPropertiesDescriptionListGroupProps = {
  customProperties: ModelRegistryCustomProperties;
  isArchive?: boolean;
  saveEditedCustomProperties: (properties: ModelRegistryCustomProperties) => Promise<unknown>;
};

const ModelPropertiesDescriptionListGroup: React.FC<ModelPropertiesDescriptionListGroupProps> = ({
  customProperties = {},
  isArchive,
  saveEditedCustomProperties,
}) => (
  <DescriptionListGroup>
    <DescriptionListDescription>
      <ModelPropertiesExpandableSection
        customProperties={customProperties}
        isArchive={isArchive}
        saveEditedCustomProperties={saveEditedCustomProperties}
        isExpandedByDefault
      />
    </DescriptionListDescription>
  </DescriptionListGroup>
);

export default ModelPropertiesDescriptionListGroup;
