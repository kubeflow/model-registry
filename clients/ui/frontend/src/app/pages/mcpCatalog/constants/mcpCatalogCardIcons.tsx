import * as React from 'react';
import { CheckCircleIcon, ClusterIcon, SecurityIcon, WrenchIcon } from '@patternfly/react-icons';

export enum McpCardIconType {
  VERIFIED_SOURCE = 'Verified source',
  SECURE_ENDPOINT = 'Secure endpoint',
  SAST = 'SAST',
  LOCAL_TO_CLUSTER = 'Local to cluster',
  READ_ONLY_TOOLS = 'Read only tools',
  REMOTE = 'Remote',
}

const GREEN_ICON_STYLE = { color: 'rgb(62, 134, 53)' };

const iconMap: Record<
  McpCardIconType,
  {
    Icon: React.ComponentType<{ className?: string; style?: React.CSSProperties }>;
    label: string;
    green?: boolean;
  }
> = {
  [McpCardIconType.VERIFIED_SOURCE]: {
    Icon: SecurityIcon,
    label: 'Verified source',
    green: true,
  },
  [McpCardIconType.SECURE_ENDPOINT]: {
    Icon: SecurityIcon,
    label: 'Secure endpoint',
    green: true,
  },
  [McpCardIconType.SAST]: {
    Icon: CheckCircleIcon,
    label: 'SAST',
    green: true,
  },
  [McpCardIconType.LOCAL_TO_CLUSTER]: {
    Icon: ClusterIcon,
    label: 'Local to cluster',
    green: false,
  },
  [McpCardIconType.READ_ONLY_TOOLS]: {
    Icon: WrenchIcon,
    label: 'Read only tools',
    green: true,
  },
  [McpCardIconType.REMOTE]: {
    Icon: ClusterIcon,
    label: 'Remote',
    green: false,
  },
};

export const getMcpCardIconConfig = (type: McpCardIconType): (typeof iconMap)[McpCardIconType] =>
  iconMap[type];

export const getMcpCardIconConfigByLabel = (
  label: string,
): {
  Icon: React.ComponentType<{ className?: string; style?: React.CSSProperties }>;
  label: string;
  green?: boolean;
} | null => {
  const entry = Object.values(iconMap).find((c) => c.label === label);
  return entry ?? null;
};

export const McpCardIcon: React.FC<{
  type: McpCardIconType;
  className?: string;
}> = ({ type, className }) => {
  const config = getMcpCardIconConfig(type);
  const { Icon, green } = config;
  return <Icon className={className} style={green ? GREEN_ICON_STYLE : undefined} />;
};

export const McpCardIconByLabel: React.FC<{
  label: string;
  className?: string;
}> = ({ label, className }) => {
  const config = getMcpCardIconConfigByLabel(label);
  if (!config) {
    return null;
  }
  const { Icon, green } = config;
  return <Icon className={className} style={green ? GREEN_ICON_STYLE : undefined} />;
};
