import * as React from 'react';
import {
  ClipboardCopy,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Flex,
  FlexItem,
  Popover,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { OutlinedQuestionCircleIcon } from '@patternfly/react-icons';
import { K8sResourceCommon } from '~/shared/types';
import '~/shared/components/NotebookController.scss';
import DashboardPopupIconButton from '~/shared/components/dashboard/DashboardPopupIconButton';

type ResourceNameTooltipProps = {
  resource: K8sResourceCommon;
  children: React.ReactNode;
  wrap?: boolean;
};

const ResourceNameTooltip: React.FC<ResourceNameTooltipProps> = ({
  children,
  resource,
  wrap = true,
}) => (
  <div style={{ display: wrap ? 'block' : 'inline-flex' }}>
    <Flex gap={{ default: 'gapXs' }} alignItems={{ default: 'alignItemsCenter' }}>
      <FlexItem>{children}</FlexItem>
      {resource.metadata?.name && (
        <Popover
          position="right"
          bodyContent={
            <Stack hasGutter>
              <StackItem>
                Resource names and types are used to find your resources in the cluster.
              </StackItem>
              <StackItem>
                <DescriptionList isCompact isHorizontal>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Resource name</DescriptionListTerm>
                    <DescriptionListDescription>
                      <ClipboardCopy
                        hoverTip="Copy"
                        clickTip="Copied"
                        variant="inline-compact"
                        data-testid="resource-name-text"
                      >
                        {resource.metadata.name}
                      </ClipboardCopy>
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Resource type</DescriptionListTerm>
                    <DescriptionListDescription data-testid="resource-kind-text">
                      {resource.kind}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                </DescriptionList>
              </StackItem>
            </Stack>
          }
        >
          <DashboardPopupIconButton
            data-testid="resource-name-icon-button"
            icon={<OutlinedQuestionCircleIcon />}
            aria-label="More info"
            style={{ paddingTop: 0, paddingBottom: 0 }}
          />
        </Popover>
      )}
    </Flex>
  </div>
);

export default ResourceNameTooltip;
