import * as React from 'react';
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Icon,
  Label,
  LabelGroup,
  PageSection,
  Popover,
  Sidebar,
  SidebarContent,
  SidebarPanel,
  Spinner,
  Alert,
  Stack,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { OutlinedClockIcon, HelpIcon } from '@patternfly/react-icons';
import { InlineTruncatedClipboardCopy } from 'mod-arch-shared';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import { CatalogArtifactList, CatalogModel } from '~/app/modelCatalogTypes';
import {
  getLabels,
  getValidatedOnPlatforms,
  getCustomPropString,
} from '~/app/pages/modelRegistry/screens/utils';
import ModelCatalogLabels from '~/app/pages/modelCatalog/components/ModelCatalogLabels';
import ExternalLink from '~/app/shared/components/ExternalLink';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import {
  getModelArtifactUri,
  hasModelArtifacts,
  isModelValidated,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CatalogModelCustomPropertyKey } from '~/concepts/modelCatalog/const';

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

  // Extract tensor type and size from customProperties
  const tensorType = model.customProperties
    ? getCustomPropString(model.customProperties, CatalogModelCustomPropertyKey.TENSOR_TYPE)
    : '';
  const size = model.customProperties
    ? getCustomPropString(model.customProperties, CatalogModelCustomPropertyKey.SIZE)
    : '';

  return (
    <PageSection hasBodyWrapper={false} isFilled padding={{ default: 'noPadding' }}>
      <Sidebar hasGutter isPanelRight>
        <SidebarContent style={{ minWidth: 0, overflow: 'hidden' }}>
          <Stack hasGutter>
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    Description
                  </Title>
                </CardHeader>
                <CardBody>
                  <Content className="pf-v6-u-text-break-word">
                    <p data-testid="model-long-description">
                      {model.description || 'No description'}
                    </p>
                  </Content>
                </CardBody>
              </Card>
            </StackItem>
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    Model card
                  </Title>
                </CardHeader>
                <CardBody>
                  {!model.readme && <p className={text.textColorDisabled}>No model card</p>}
                  {model.readme && (
                    <MarkdownComponent
                      data={model.readme}
                      dataTestId="model-card-markdown"
                      maxHeading={3}
                    />
                  )}
                </CardBody>
              </Card>
            </StackItem>
          </Stack>
        </SidebarContent>
        <SidebarPanel width={{ default: 'width_33' }}>
          <Card>
            <CardHeader>
              <Title headingLevel="h2" size="lg">
                Model details
              </Title>
            </CardHeader>
            <CardBody>
              <DescriptionList>
                <DescriptionListGroup>
                  <DescriptionListTerm>Labels</DescriptionListTerm>
                  <DescriptionListDescription>
                    <ModelCatalogLabels
                      tasks={model.tasks ?? []}
                      labels={allLabels.filter((label) => label !== 'validated')}
                      numLabels={isValidated ? 2 : 3}
                    />
                  </DescriptionListDescription>
                </DescriptionListGroup>
                <DescriptionListGroup>
                  <DescriptionListTerm>Tensor type</DescriptionListTerm>
                  <DescriptionListDescription>{tensorType || 'N/A'}</DescriptionListDescription>
                </DescriptionListGroup>
                <DescriptionListGroup>
                  <DescriptionListTerm>Size</DescriptionListTerm>
                  <DescriptionListDescription>{size || 'N/A'}</DescriptionListDescription>
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
                    <DescriptionListTerm>
                      Validated on{' '}
                      <Popover bodyContent="The platform versions where this model has been validated for operational readiness and compatibility.">
                        <Button
                          variant="plain"
                          aria-label="More info for validated on"
                          className="pf-v6-u-p-xs"
                          icon={<HelpIcon />}
                        />
                      </Popover>
                    </DescriptionListTerm>
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
            </CardBody>
          </Card>
        </SidebarPanel>
      </Sidebar>
    </PageSection>
  );
};

export default ModelDetailsView;
