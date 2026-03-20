import * as React from 'react';
import {
  CodeBlock,
  CodeBlockCode,
  DrawerActions,
  DrawerCloseButton,
  DrawerHead,
  DrawerPanelBody,
} from '@patternfly/react-core';
import sampleCatalogYamlContent from '~/app/pages/modelCatalogSettings/sample-catalog.yaml';
import { EXPECTED_YAML_FORMAT_LABEL } from '~/app/pages/modelCatalogSettings/constants';

type ExpectedYamlFormatDrawerPanelProps = {
  onClose: () => void;
};

export const ExpectedYamlFormatDrawerPanel: React.FC<ExpectedYamlFormatDrawerPanelProps> = ({
  onClose,
}) => (
  <>
    <DrawerHead>
      <span data-testid="expected-format-drawer-title">{EXPECTED_YAML_FORMAT_LABEL}</span>
      <DrawerActions>
        <DrawerCloseButton
          onClose={onClose}
          aria-label="Close drawer"
          data-testid="expected-format-drawer-close"
        />
      </DrawerActions>
    </DrawerHead>
    <DrawerPanelBody hasNoPadding>
      <CodeBlock>
        <CodeBlockCode>{sampleCatalogYamlContent}</CodeBlockCode>
      </CodeBlock>
    </DrawerPanelBody>
  </>
);
