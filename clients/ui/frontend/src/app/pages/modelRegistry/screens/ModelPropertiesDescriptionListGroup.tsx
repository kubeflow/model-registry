import * as React from 'react';
import { DashboardDescriptionListGroup } from 'mod-arch-shared';
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
  <DashboardDescriptionListGroup
    title=""
    isEmpty={Object.keys(customProperties).length === 0}
    contentWhenEmpty="No properties"
  >
    <ModelPropertiesExpandableSection
      customProperties={customProperties}
      isArchive={isArchive}
      saveEditedCustomProperties={saveEditedCustomProperties}
      isExpandedByDefault
    />
  </DashboardDescriptionListGroup>
);

export default ModelPropertiesDescriptionListGroup;
