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
  Spinner,
  Alert,
} from '@patternfly/react-core';
import { OutlinedClockIcon } from '@patternfly/react-icons';
import { InlineTruncatedClipboardCopy } from 'mod-arch-shared';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import { CatalogModel, CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { useCatalogModelArtifacts } from '~/app/hooks/modelCatalog/useCatalogModelArtifacts';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import ModelCatalogLabels from '~/app/pages/modelCatalog/components/ModelCatalogLabels';
import ExternalLink from '~/app/shared/components/ExternalLink';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';

type ModelDetailsViewProps = {
  model: CatalogModel;
  decodedParams: CatalogModelDetailsParams;
};

const ModelDetailsView: React.FC<ModelDetailsViewProps> = ({ model, decodedParams }) => {
  const [artifacts, artifactLoaded, artifactsLoadError] = useCatalogModelArtifacts(
    decodedParams.sourceId || '',
    encodeURIComponent(`${decodedParams.modelName}`),
  );

  // Extract all labels from customProperties
  const allLabels = model.customProperties ? getLabels(model.customProperties) : [];

  return (
    <PageSection hasBodyWrapper={false} isFilled>
      <Sidebar hasBorder hasGutter isPanelRight>
        <SidebarContent>
          <Content>
            <h2>Description</h2>
            <p data-testid="model-long-description">{model.description || 'No description'}</p>
            <h2>Model card</h2>
            {!model.readme && <p className={text.textColorDisabled}>No model card</p>}
          </Content>
          {model.readme && (
            <MarkdownComponent
              data={model.readme}
              dataTestId="model-card-markdown"
              maxHeading={3}
            />
          )}
        </SidebarContent>
        <SidebarPanel>
          <DescriptionList isFillColumns>
            <DescriptionListGroup>
              <DescriptionListTerm>Labels</DescriptionListTerm>
              <DescriptionListDescription>
                <ModelCatalogLabels
                  tasks={model.tasks ?? []}
                  license={model.license}
                  labels={allLabels}
                />
              </DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>License</DescriptionListTerm>
              <ExternalLink
                text="Agreement"
                to={model.licenseLink || ''}
                testId="model-license-link"
              />
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Provider</DescriptionListTerm>
              <DescriptionListDescription>{model.provider || 'N/A'}</DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Model location</DescriptionListTerm>
              {artifactsLoadError ? (
                <Alert variant="danger" isInline title={artifactsLoadError.name}>
                  {artifactsLoadError.message}
                </Alert>
              ) : !artifactLoaded ? (
                <Spinner size="sm" />
              ) : artifacts.items.length > 0 ? (
                <InlineTruncatedClipboardCopy
                  testId="source-image-location"
                  textToCopy={artifacts.items.map((artifact) => artifact.uri)[0] || ''}
                />
              ) : (
                <p className={text.textColorDisabled}>No artifacts available</p>
              )}
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Last modified</DescriptionListTerm>
              <DescriptionListDescription>
                <Icon isInline style={{ marginRight: 4 }}>
                  <OutlinedClockIcon />
                </Icon>
                <ModelTimestamp timeSinceEpoch={model.lastUpdateTimeSinceEpoch} />
              </DescriptionListDescription>
            </DescriptionListGroup>
            <DescriptionListGroup>
              <DescriptionListTerm>Published</DescriptionListTerm>
              <DescriptionListDescription>
                <Icon isInline style={{ marginRight: 4 }}>
                  <OutlinedClockIcon />
                </Icon>
                <ModelTimestamp timeSinceEpoch={model.createTimeSinceEpoch} />
              </DescriptionListDescription>
            </DescriptionListGroup>
          </DescriptionList>
        </SidebarPanel>
      </Sidebar>
    </PageSection>
  );
};

export default ModelDetailsView;
