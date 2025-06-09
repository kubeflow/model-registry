import * as React from 'react';
import { Stack, StackItem } from '@patternfly/react-core';
import { FieldGroupHelpLabelIcon } from 'mod-arch-shared';

const ResourceNameDefinitionTooltip: React.FC = () => (
  <FieldGroupHelpLabelIcon
    content={
      <Stack hasGutter>
        <StackItem>
          The resource name is used to identify your resource, and is generated based on the name
          you enter.
        </StackItem>
        <StackItem>The resource name cannot be edited after creation.</StackItem>
      </Stack>
    }
  />
);

export default ResourceNameDefinitionTooltip;
