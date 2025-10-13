import * as React from 'react';
import {
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Icon,
  Label,
  LabelGroup,
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
import { CatalogArtifactList, CatalogModel } from '~/app/modelCatalogTypes';
import { getLabels, getValidatedOnPlatforms } from '~/app/pages/modelRegistry/screens/utils';
import ModelCatalogLabels from '~/app/pages/modelCatalog/components/ModelCatalogLabels';
import ExternalLink from '~/app/shared/components/ExternalLink';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import {
  getModelArtifactUri,
  hasModelArtifacts,
  isModelValidated,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type ModelDetailsViewProps = {
  model: CatalogModel;
  artifacts: CatalogArtifactList;
  artifactLoaded: boolean;
  artifactsLoadError: Error | undefined;
};

const ModelDetailsView: React.FC<ModelDetailsViewProps> = ({
  model,
  artifacts,
  artifactLoaded,
  artifactsLoadError,
}) => {
  // Extract all labels from customProperties
  const allLabels = model.customProperties ? getLabels(model.customProperties) : [];
  const isValidated = isModelValidated(model);

  // Extract validated_on platforms
  const validatedOnPlatforms = getValidatedOnPlatforms(model.customProperties);

  return (
    <PageSection hasBodyWrapper={false} isFilled>
      <Sidebar hasBorder hasGutter isPanelRight>
        <SidebarContent style={{ minWidth: 0, overflow: 'hidden' }}>
          <Content style={{ wordBreak: 'break-word' }}>
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
        <SidebarPanel width={{ default: 'width_33' }} style={{ minWidth: 0 }}>
          <DescriptionList>
            <DescriptionListGroup>
              <DescriptionListTerm>Labels</DescriptionListTerm>
              <DescriptionListDescription>
                <ModelCatalogLabels
                  tasks={model.tasks ?? []}
                  license={model.license}
                  labels={allLabels}
                  numLabels={isValidated ? 2 : 3}
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
            {validatedOnPlatforms.length > 0 && (
              <DescriptionListGroup>
                <DescriptionListTerm>Validated on</DescriptionListTerm>
                <DescriptionListDescription>
                  <LabelGroup numLabels={5} isCompact>
                    {validatedOnPlatforms.map((platform) => (
                      <Label data-testid="validated-on-label" key={platform} variant="outline">
                        {platform}
                      </Label>
                    ))}
                  </LabelGroup>
                </DescriptionListDescription>
              </DescriptionListGroup>
            )}
            <DescriptionListGroup>
              <DescriptionListTerm>Model location</DescriptionListTerm>
              {artifactsLoadError ? (
                <Alert variant="danger" isInline title={artifactsLoadError.name}>
                  {artifactsLoadError.message}
                </Alert>
              ) : !artifactLoaded ? (
                <Spinner size="sm" />
              ) : artifacts.items.length > 0 && hasModelArtifacts(artifacts.items) ? (
                <DescriptionListDescription>
                  <InlineTruncatedClipboardCopy
                    testId="source-image-location"
                    textToCopy={getModelArtifactUri(artifacts.items) || ''}
                  />
                </DescriptionListDescription>
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
