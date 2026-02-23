import * as React from 'react';
import { createPortal } from 'react-dom';
import { Button, CodeBlock, CodeBlockCode } from '@patternfly/react-core';
import TimesIcon from '@patternfly/react-icons/dist/esm/icons/times-icon';
import { EXPECTED_YAML_FORMAT_CONTENT } from '~/app/pages/modelCatalogSettings/expectedYamlFormatContent';

const DRAWER_TITLE = 'View expected file format';
const PRIMARY_APP_CONTAINER_ID = 'primary-app-container';
const DRAWER_BOTTOM_OFFSET_PX = 85;

type ContainerRect = {
  top: number;
  left: number;
  width: number;
  height: number;
  drawerHeight: number;
};

function useDrawerBounds(isOpen: boolean) {
  const [bounds, setBounds] = React.useState<ContainerRect | null>(null);

  React.useLayoutEffect(() => {
    if (!isOpen) {
      setBounds(null);
      return;
    }

    const updateBounds = () => {
      const el = document.getElementById(PRIMARY_APP_CONTAINER_ID);
      if (!el) return;
      const rect = el.getBoundingClientRect();
      const drawerHeight = Math.max(0, rect.height - DRAWER_BOTTOM_OFFSET_PX);
      setBounds({
        top: rect.top,
        left: rect.left,
        width: rect.width,
        height: rect.height,
        drawerHeight,
      });
    };

    updateBounds();
    const container = document.getElementById(PRIMARY_APP_CONTAINER_ID);
    if (!container) return;

    const resizeObserver = new ResizeObserver(updateBounds);
    resizeObserver.observe(container);
    container.addEventListener('scroll', updateBounds, true);

    return () => {
      resizeObserver.disconnect();
      container.removeEventListener('scroll', updateBounds, true);
    };
  }, [isOpen]);

  return bounds;
}

type ExpectedYamlFormatDrawerProps = {
  isOpen: boolean;
  onClose: () => void;
};

const ExpectedYamlFormatDrawer: React.FC<ExpectedYamlFormatDrawerProps> = ({ isOpen, onClose }) => {
  const drawerBounds = useDrawerBounds(isOpen);

  const overlay =
    isOpen &&
    drawerBounds && (
      <div
        role="presentation"
        style={{
          position: 'fixed',
          top: drawerBounds.top,
          left: drawerBounds.left,
          width: drawerBounds.width,
          height: drawerBounds.height,
          zIndex: 100,
          pointerEvents: 'none',
          overflow: 'hidden',
        }}
      >
        <div
          role="region"
          aria-label={DRAWER_TITLE}
          style={{
            position: 'absolute',
            right: 0,
            top: 0,
            width: '51%',
            height: `${drawerBounds.drawerHeight}px`,
            maxHeight: `${drawerBounds.drawerHeight}px`,
            pointerEvents: 'auto',
            display: 'flex',
            flexDirection: 'column',
            minHeight: 0,
            backgroundColor: '#fff',
            borderLeft: '1px solid var(--pf-v6-global--BorderColor--100)',
            boxShadow: '-4px 0 8px rgba(0, 0, 0, 0.1)',
          }}
        >
          <div
            className="pf-v6-c-drawer__head pf-v6-c-page__main-breadcrumb"
            style={{
              display: 'flex',
              flexDirection: 'row',
              flexWrap: 'nowrap',
              alignItems: 'center',
              justifyContent: 'space-between',
              gap: 'var(--pf-v6-global--spacer--md)',
            padding: '16px',
            flexShrink: 0,
            borderBottom: '1px solid var(--pf-v6-global--BorderColor--100)',
            minHeight: 'var(--pf-v6-c-breadcrumb__item--FontSize)',
            }}
          >
            <h2
              data-testid="expected-format-drawer-title"
              style={{
                minWidth: 0,
                flex: '1 1 auto',
                margin: 0,
                fontSize: 'large',
                fontWeight: 500,
              }}
            >
              {DRAWER_TITLE}
            </h2>
            <span style={{ display: 'inline-flex', flexShrink: 0, alignItems: 'center' }}>
              <Button
                variant="plain"
                aria-label="Close drawer"
                data-testid="expected-format-drawer-close"
                icon={<TimesIcon />}
                onClick={onClose}
                className="pf-v6-c-drawer__close"
              />
            </span>
          </div>
          <div
            style={{
              flex: '1 1 auto',
              minHeight: 0,
              overflow: 'auto',
              padding: '16px',
            }}
          >
            <CodeBlock>
              <CodeBlockCode>{EXPECTED_YAML_FORMAT_CONTENT}</CodeBlockCode>
            </CodeBlock>
          </div>
        </div>
      </div>
    );

  const mainContainer =
    typeof document !== 'undefined' ? document.getElementById(PRIMARY_APP_CONTAINER_ID) : null;

  if (!isOpen || !mainContainer || !overlay) return null;

  return createPortal(overlay, mainContainer);
};

export default ExpectedYamlFormatDrawer;
