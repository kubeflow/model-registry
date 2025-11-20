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
} from '@patternfly/react-core';
import { CubesIcon } from '@patternfly/react-icons';
import { PAGE_TITLES } from '~/app/pages/modelCatalogSettings/constants';
import PreviewButton from './PreviewButton';

type PreviewPanelProps = {
  isPreviewEnabled: boolean;
  onPreview: () => void;
};

const PreviewPanel: React.FC<PreviewPanelProps> = ({ isPreviewEnabled, onPreview }) => (
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
          variant="secondary"
          testId="preview-button-header"
        />
      </FlexItem>
    </Flex>
    <EmptyState
      icon={CubesIcon}
      titleText={PAGE_TITLES.PREVIEW_MODELS}
      variant={EmptyStateVariant.sm}
    >
      <EmptyStateBody>
        To view the models from this source that will appear in the model catalog with your current
        configuration, complete all required fields, then click <strong>Preview</strong>.
      </EmptyStateBody>
      <EmptyStateFooter>
        <EmptyStateActions>
          <PreviewButton
            onClick={onPreview}
            isDisabled={!isPreviewEnabled}
            variant="link"
            testId="preview-button-panel"
          />
        </EmptyStateActions>
      </EmptyStateFooter>
    </EmptyState>
  </div>
);

export default PreviewPanel;
