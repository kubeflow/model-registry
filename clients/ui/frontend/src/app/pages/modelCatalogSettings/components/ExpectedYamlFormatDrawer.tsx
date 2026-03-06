import * as React from 'react';
import { createPortal } from 'react-dom';
import {
  CodeBlock,
  CodeBlockCode,
  DrawerActions,
  DrawerCloseButton,
  DrawerHead,
  DrawerPanelBody,
} from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import sampleCatalogYamlContent from '~/app/pages/modelCatalogSettings/sample-catalog.yaml';
import {
  EXPECTED_YAML_FORMAT_LABEL,
  PRIMARY_APP_CONTAINER_ID,
} from '~/app/pages/modelCatalogSettings/constants';

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
    <DrawerPanelBody style={{ flex: 1, minHeight: 0, overflow: 'auto' }}>
      <CodeBlock>
        <CodeBlockCode>{sampleCatalogYamlContent}</CodeBlockCode>
      </CodeBlock>
    </DrawerPanelBody>
  </>
);

const panelStyle: React.CSSProperties = {
  position: 'absolute',
  top: 0,
  right: 0,
  bottom: 0,
  width: '50%',
  zIndex: 400,
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
  padding: 16,
  backgroundColor: 'var(--pf-t--global--background--color--primary--default)',
  boxShadow: 'var(--pf-t--global--box-shadow--lg--left)',
};

type ExpectedYamlFormatDrawerProps = {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
};

const ExpectedYamlFormatDrawer: React.FC<ExpectedYamlFormatDrawerProps> = ({
  isOpen,
  onClose,
  children,
}) => {
  const { isMUITheme } = useThemeContext();
  const container =
    typeof document !== 'undefined' ? document.getElementById(PRIMARY_APP_CONTAINER_ID) : null;

  React.useEffect(() => {
    if (!isOpen || !container) {
      return;
    }
    const prev = container.style.position;
    if (!prev || prev === 'static') {
      container.style.position = 'relative';
    }
    return () => {
      container.style.position = prev;
    };
  }, [isOpen, container]);

  const panel =
    isOpen && container
      ? createPortal(
          <div
            role="region"
            aria-label={EXPECTED_YAML_FORMAT_LABEL}
            style={{
              ...panelStyle,
              ...(isMUITheme && { maxWidth: 'none' }),
            }}
          >
            <ExpectedYamlFormatDrawerPanel onClose={onClose} />
          </div>,
          container,
        )
      : null;

  return (
    <>
      {children}
      {panel}
    </>
  );
};

export default ExpectedYamlFormatDrawer;
