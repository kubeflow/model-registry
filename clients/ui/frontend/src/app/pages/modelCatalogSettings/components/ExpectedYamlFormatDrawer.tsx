import * as React from 'react';
import { createPortal } from 'react-dom';
import {
  CodeBlock,
  CodeBlockCode,
  Drawer,
  DrawerActions,
  DrawerCloseButton,
  DrawerContent,
  DrawerHead,
  DrawerPanelBody,
  DrawerPanelContent,
} from '@patternfly/react-core';
import sampleCatalogYamlContent from '@repo/manifests/kustomize/options/catalog/base/sample-catalog.yaml';

export const EXPECTED_FORMAT_DRAWER_TITLE = 'View expected file format';

const PRIMARY_APP_CONTAINER_ID = 'primary-app-container';

type ExpectedYamlFormatDrawerPanelProps = {
  onClose: () => void;
};

export const ExpectedYamlFormatDrawerPanel: React.FC<ExpectedYamlFormatDrawerPanelProps> = ({
  onClose,
}) => (
  <DrawerPanelContent
    widths={{ default: 'width_50' }}
    role="region"
    aria-label={EXPECTED_FORMAT_DRAWER_TITLE}
    style={{ display: 'flex', flexDirection: 'column', minHeight: 0 }}
  >
    <DrawerHead>
      <span data-testid="expected-format-drawer-title">{EXPECTED_FORMAT_DRAWER_TITLE}</span>
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
  </DrawerPanelContent>
);

type ExpectedYamlFormatDrawerProps = {
  isOpen: boolean;
  onClose: () => void;
};

type ContainerBounds = { top: number; left: number; width: number; height: number };

function useContainerBounds(
  container: HTMLElement | null,
  isActive: boolean,
): ContainerBounds | null {
  const [bounds, setBounds] = React.useState<ContainerBounds | null>(null);

  React.useLayoutEffect(() => {
    if (!isActive || !container) {
      setBounds(null);
      return;
    }
    const update = () => {
      const rect = container.getBoundingClientRect();
      setBounds({
        top: rect.top,
        left: rect.left,
        width: rect.width,
        height: rect.height,
      });
    };
    update();
    const ro = new ResizeObserver(update);
    ro.observe(container);
    window.addEventListener('scroll', update, true);
    return () => {
      ro.disconnect();
      window.removeEventListener('scroll', update, true);
    };
  }, [isActive, container]);

  return bounds;
}

const ExpectedYamlFormatDrawer: React.FC<ExpectedYamlFormatDrawerProps> = ({ isOpen, onClose }) => {
  const container =
    typeof document !== 'undefined' ? document.getElementById(PRIMARY_APP_CONTAINER_ID) : null;
  const bounds = useContainerBounds(container, isOpen);

  if (!isOpen || !container || !bounds) {
    return null;
  }

  const drawerNode = (
    <div
      style={{
        position: 'fixed',
        top: bounds.top,
        left: bounds.left,
        width: bounds.width,
        height: bounds.height,
        zIndex: 100,
        pointerEvents: 'auto',
        overflow: 'hidden',
        display: 'flex',
        flexDirection: 'column',
        boxSizing: 'border-box',
      }}
    >
      <Drawer
        isExpanded
        position="end"
        style={{
          flex: 1,
          minHeight: 0,
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <DrawerContent panelContent={<ExpectedYamlFormatDrawerPanel onClose={onClose} />}>
          <div />
        </DrawerContent>
      </Drawer>
    </div>
  );

  return createPortal(drawerNode, document.body);
};

export default ExpectedYamlFormatDrawer;
