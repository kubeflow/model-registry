import * as React from 'react';
import { Flex, FlexItem, Icon, Tooltip } from '@patternfly/react-core';
import {
  CheckCircleIcon,
  ExclamationCircleIcon,
  ShieldAltIcon,
  LockIcon,
  CodeIcon,
  EyeIcon,
} from '@patternfly/react-icons';
import { McpSecurityIndicator } from '~/app/pages/mcpCatalog/types';

type McpSecurityIndicatorsProps = {
  indicators: McpSecurityIndicator;
};

type IndicatorItemProps = {
  isEnabled: boolean;
  enabledLabel: string;
  disabledLabel: string;
  icon: React.ReactNode;
};

const IndicatorItem: React.FC<IndicatorItemProps> = ({
  isEnabled,
  enabledLabel,
  disabledLabel,
  icon,
}) => (
  <FlexItem>
    <Tooltip content={isEnabled ? enabledLabel : disabledLabel}>
      <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapXs' }}>
        <Icon status={isEnabled ? 'success' : 'warning'}>
          {isEnabled ? <CheckCircleIcon /> : <ExclamationCircleIcon />}
        </Icon>
        <span style={{ fontSize: 'var(--pf-t--global--font--size--sm)' }}>{icon}</span>
        <span style={{ fontSize: 'var(--pf-t--global--font--size--sm)' }}>
          {isEnabled ? enabledLabel : disabledLabel}
        </span>
      </Flex>
    </Tooltip>
  </FlexItem>
);

const McpSecurityIndicators: React.FC<McpSecurityIndicatorsProps> = ({ indicators }) => (
  <Flex direction={{ default: 'column' }} gap={{ default: 'gapXs' }}>
    <IndicatorItem
      isEnabled={indicators.verifiedSource}
      enabledLabel="Verified source"
      disabledLabel="Unverified source"
      icon={<ShieldAltIcon />}
    />
    <IndicatorItem
      isEnabled={indicators.secureEndpoint}
      enabledLabel="Secure endpoint"
      disabledLabel="Insecure endpoint"
      icon={<LockIcon />}
    />
    {indicators.sast && (
      <IndicatorItem
        isEnabled={indicators.sast}
        enabledLabel="SAST"
        disabledLabel="No SAST"
        icon={<CodeIcon />}
      />
    )}
    {indicators.readOnlyTools && (
      <IndicatorItem
        isEnabled={indicators.readOnlyTools}
        enabledLabel="Read only tools"
        disabledLabel="Read/Write tools"
        icon={<EyeIcon />}
      />
    )}
  </Flex>
);

export default McpSecurityIndicators;
