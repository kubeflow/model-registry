import * as React from 'react';
import {
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Icon,
  PageSection,
  Sidebar,
  SidebarContent,
  SidebarPanel,
  Label,
} from '@patternfly/react-core';
import { TagIcon, OutlinedClockIcon } from '@patternfly/react-icons';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import {
  extractVersionTag,
  filterNonVersionTags,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogLabels from '~/app/pages/modelCatalog/components/ModelCatalogLabels';

type ModelDetailsViewProps = {
  model: ModelCatalogItem;
};

const ModelDetailsView: React.FC<ModelDetailsViewProps> = ({ model }) => {
  const versionTag = extractVersionTag(model.tags);
  const nonVersionTags = filterNonVersionTags(model.tags) ?? [];

  return (
    <PageSection hasBodyWrapper={false} isFilled>
      <Sidebar hasBorder hasGutter isPanelRight>
        <SidebarContent>
          <Content>
            <h2>Description</h2>
            <p data-testid="model-long-description">{model.description || 'No description'}</p>
            <h2>Model card</h2>
            <p className="pf-v5-u-color-200">No model card</p>
          </Content>
        </SidebarContent>
        <SidebarPanel>
          <DescriptionList isFillColumns>
            <DescriptionListGroup>
              <DescriptionListTerm>Version</DescriptionListTerm>
              <DescriptionListDescription>
                <Label variant="outline" icon={<TagIcon />}>
                  {versionTag || 'N/A'}
                </Label>
              </DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Labels</DescriptionListTerm>
              <DescriptionListDescription>
                <ModelCatalogLabels
                  tags={nonVersionTags}
                  framework={model.framework}
                  task={model.task}
                  license={model.license}
                />
              </DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>License</DescriptionListTerm>
              <DescriptionListDescription>{model.license || 'N/A'}</DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Provider</DescriptionListTerm>
              <DescriptionListDescription>{model.provider || 'N/A'}</DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Model location</DescriptionListTerm>
              <DescriptionListDescription>{model.url || 'N/A'}</DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Last modified</DescriptionListTerm>
              <DescriptionListDescription>
                <Icon isInline style={{ marginRight: 4 }}>
                  <OutlinedClockIcon />
                </Icon>
                {model.updatedAt || 'N/A'}
              </DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Published</DescriptionListTerm>
              <DescriptionListDescription>
                <Icon isInline style={{ marginRight: 4 }}>
                  <OutlinedClockIcon />
                </Icon>
                {model.createdAt || 'N/A'}
              </DescriptionListDescription>
            </DescriptionListGroup>
          </DescriptionList>
        </SidebarPanel>
      </Sidebar>
    </PageSection>
  );
};

export default ModelDetailsView;
