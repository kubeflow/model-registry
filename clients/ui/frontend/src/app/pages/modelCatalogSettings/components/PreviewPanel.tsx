import * as React from 'react';
import {
  EmptyState,
  EmptyStateVariant,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateActions,
  Flex,
  FlexItem,
  Title,
  Tabs,
  Tab,
  TabTitleText,
  Alert,
  List,
  ListItem,
  Spinner,
  Pagination,
  PaginationVariant,
  AlertActionLink,
} from '@patternfly/react-core';
import { CubesIcon, CheckCircleIcon, TimesCircleIcon } from '@patternfly/react-icons';
import { PAGE_TITLES } from '~/app/pages/modelCatalogSettings/constants';
import { CatalogSourcePreviewResult } from '~/app/modelCatalogTypes';
import PreviewButton from './PreviewButton';

type PreviewPanelProps = {
  isPreviewEnabled: boolean;
  isLoading: boolean;
  onPreview: () => void;
  previewResult?: CatalogSourcePreviewResult;
  previewError?: Error;
  hasFormChanged: boolean;
};

const PreviewPanel: React.FC<PreviewPanelProps> = ({
  isPreviewEnabled,
  isLoading,
  onPreview,
  previewResult,
  previewError,
  hasFormChanged,
}) => {
  const [activeTabKey, setActiveTabKey] = React.useState<string | number>(0);
  const [page, setPage] = React.useState(1);
  const [perPage, setPerPage] = React.useState(10);

  const handleTabSelect = (_event: React.MouseEvent, tabIndex: string | number) => {
    setActiveTabKey(tabIndex);
    setPage(1); // Reset to first page when switching tabs
  };

  const filteredItems = React.useMemo(() => {
    if (!previewResult) {
      return [];
    }
    if (activeTabKey === 0) {
      return previewResult.items.filter((item) => item.included);
    }
    return previewResult.items.filter((item) => !item.included);
  }, [previewResult, activeTabKey]);

  const paginatedItems = React.useMemo(() => {
    const startIdx = (page - 1) * perPage;
    const endIdx = startIdx + perPage;
    return filteredItems.slice(startIdx, endIdx);
  }, [filteredItems, page, perPage]);

  const onSetPage = (
    _event: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPage: number,
  ) => {
    setPage(newPage);
  };

  const onPerPageSelect = (
    _event: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
  ) => {
    setPerPage(newPerPage);
    setPage(1);
  };

  const renderEmptyState = () => {
    if (previewError) {
      return (
        <EmptyState
          icon={TimesCircleIcon}
          titleText="Failed to preview the results"
          variant={EmptyStateVariant.sm}
        >
          <EmptyStateBody>{previewError.message}</EmptyStateBody>
          <EmptyStateFooter>
            <EmptyStateActions>
              <PreviewButton
                onClick={onPreview}
                isDisabled={!isPreviewEnabled}
                isLoading={isLoading}
                variant="link"
                testId="preview-button-panel-retry"
              />
            </EmptyStateActions>
          </EmptyStateFooter>
        </EmptyState>
      );
    }

    return (
      <EmptyState
        icon={CubesIcon}
        titleText={PAGE_TITLES.PREVIEW_MODELS}
        variant={EmptyStateVariant.sm}
      >
        <EmptyStateBody>
          To view the models from this source that will appear in the model catalog with your
          current configuration, complete all required fields, then click <strong>Preview</strong>.
        </EmptyStateBody>
        <EmptyStateFooter>
          <EmptyStateActions>
            <PreviewButton
              onClick={onPreview}
              isDisabled={!isPreviewEnabled}
              isLoading={isLoading}
              variant="link"
              testId="preview-button-panel"
            />
          </EmptyStateActions>
        </EmptyStateFooter>
      </EmptyState>
    );
  };

  const renderContent = () => {
    if (isLoading) {
      return (
        <div className="pf-v6-u-text-align-center pf-v6-u-py-xl">
          <Spinner size="xl" aria-label="Loading preview" />
        </div>
      );
    }

    if (!previewResult || previewError) {
      return renderEmptyState();
    }

    return (
      <>
        <Tabs activeKey={activeTabKey} onSelect={handleTabSelect} aria-label="Preview tabs">
          <Tab eventKey={0} title={<TabTitleText>Models included</TabTitleText>} />
          <Tab eventKey={1} title={<TabTitleText>Models excluded</TabTitleText>} />
        </Tabs>
        <div className="pf-v6-u-mt-md">
          {hasFormChanged && (
            <Alert
              variant="info"
              isInline
              title="The preview needs to be refreshed after any changes are made"
              className="pf-v6-u-mb-md"
              actionLinks={
                <AlertActionLink onClick={onPreview} data-testid="refresh-preview-link">
                  Refresh the preview
                </AlertActionLink>
              }
            />
          )}
          {paginatedItems.length > 0 ? (
            <>
              <Flex
                justifyContent={{ default: 'justifyContentSpaceBetween' }}
                alignItems={{ default: 'alignItemsCenter' }}
              >
                <FlexItem>
                  <strong>
                    {activeTabKey === 0
                      ? `${previewResult.summary.includedModels} of ${previewResult.summary.totalModels} models included:`
                      : `${previewResult.summary.excludedModels} of ${previewResult.summary.totalModels} models excluded:`}
                  </strong>
                </FlexItem>
                <FlexItem>
                  <Pagination
                    itemCount={filteredItems.length}
                    perPage={perPage}
                    page={page}
                    onSetPage={onSetPage}
                    onPerPageSelect={onPerPageSelect}
                    variant={PaginationVariant.top}
                  />
                </FlexItem>
              </Flex>
              <List isPlain className="pf-v6-u-mt-md">
                {paginatedItems.map((model) => (
                  <ListItem
                    key={model.name}
                    icon={
                      model.included ? (
                        <CheckCircleIcon color="green" />
                      ) : (
                        <TimesCircleIcon color="red" />
                      )
                    }
                  >
                    {model.name}
                  </ListItem>
                ))}
              </List>
              <Pagination
                itemCount={filteredItems.length}
                perPage={perPage}
                page={page}
                onSetPage={onSetPage}
                onPerPageSelect={onPerPageSelect}
                variant={PaginationVariant.bottom}
              />
            </>
          ) : (
            <EmptyState
              variant={EmptyStateVariant.sm}
              titleText={`No models ${activeTabKey === 0 ? 'included' : 'excluded'}`}
            >
              <EmptyStateBody>
                No models from this source are visible in the model catalog
              </EmptyStateBody>
            </EmptyState>
          )}
        </div>
      </>
    );
  };

  return (
    <div data-testid="preview-panel" className="pf-v6-u-h-100">
      <Flex
        justifyContent={{ default: 'justifyContentSpaceBetween' }}
        alignItems={{ default: 'alignItemsCenter' }}
        className="pf-v6-u-mb-md"
      >
        <FlexItem>
          <Title headingLevel="h2" size="lg">
            {PAGE_TITLES.MODEL_CATALOG_PREVIEW}
          </Title>
        </FlexItem>
        <FlexItem>
          <PreviewButton
            onClick={onPreview}
            isDisabled={!isPreviewEnabled}
            isLoading={isLoading}
            variant="secondary"
            testId="preview-button-header"
          />
        </FlexItem>
      </Flex>
      {renderContent()}
    </div>
  );
};

export default PreviewPanel;
