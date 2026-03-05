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
  <DrawerPanelContent
    widths={{ default: 'width_100' }}
    role="region"
    aria-label={EXPECTED_YAML_FORMAT_LABEL}
    style={{
      display: 'flex',
      flexDirection: 'column',
      minHeight: 0,
      width: '100%',
      maxWidth: 'none',
      flexBasis: '100%',
    }}
  >
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
    <DrawerPanelBody style={{ flex: 1, minHeight: 200, overflow: 'auto' }}>
      <CodeBlock>
        <CodeBlockCode>{sampleCatalogYamlContent}</CodeBlockCode>
      </CodeBlock>
    </DrawerPanelBody>
  </DrawerPanelContent>
);

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
  const wrapperRef = React.useRef<HTMLDivElement>(null);
  const container =
    typeof document !== 'undefined' ? document.getElementById(PRIMARY_APP_CONTAINER_ID) : null;

  React.useEffect(() => {
    if (isOpen && !container) {
      // eslint-disable-next-line no-console
      console.warn(
        `ExpectedYamlFormatDrawer: container #${PRIMARY_APP_CONTAINER_ID} not found. The drawer will not render.`,
      );
    }
  }, [isOpen, container]);

  const bounds = useContainerBounds(container, isOpen);

  React.useEffect(() => {
    const el = wrapperRef.current;
    if (!el) {
      return;
    }
    if (isMUITheme) {
      el.style.setProperty('--mui-drawer__panel--MaxWidth', 'none');
      const applyPanelMaxWidth = () => {
        const panel = el.querySelector('.pf-v6-c-drawer__panel');
        if (panel instanceof HTMLElement) {
          panel.style.maxWidth = 'none';
        }
      };
      applyPanelMaxWidth();
      const id = window.setTimeout(applyPanelMaxWidth, 0);
      const observer = new MutationObserver(applyPanelMaxWidth);
      observer.observe(el, { childList: true, subtree: true });
      return () => {
        window.clearTimeout(id);
        observer.disconnect();
      };
    }
    el.style.removeProperty('--mui-drawer__panel--MaxWidth');
    return undefined;
  }, [isMUITheme, isOpen]);

  if (!isOpen) {
    return <>{children}</>;
  }

  if (!bounds) {
    return <>{children}</>;
  }

  const halfWidth = bounds.width / 2;
  const overlay = (
    <div
      ref={wrapperRef}
      style={{
        position: 'fixed',
        top: bounds.top,
        left: bounds.left + halfWidth,
        width: halfWidth,
        height: bounds.height,
        zIndex: 100,
        pointerEvents: 'auto',
        overflow: 'hidden',
        display: 'flex',
        flexDirection: 'column',
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

  return (
    <>
      {children}
      {createPortal(overlay, document.body)}
    </>
  );
};

export default ExpectedYamlFormatDrawer;
