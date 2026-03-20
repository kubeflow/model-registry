import * as React from 'react';
import { CheckCircleIcon, ServerIcon, ShieldAltIcon, ToolsIcon } from '@patternfly/react-icons';

export enum McpCardIconType {
  VERIFIED_SOURCE = 'Verified source',
  SECURE_ENDPOINT = 'Secure endpoint',
  SAST = 'SAST',
  LOCAL_TO_CLUSTER = 'Local to cluster',
  READ_ONLY_TOOLS = 'Read only tools',
  RED_HAT_PARTNER = 'Red Hat partner',
  REMOTE = 'Remote',
}

const GREEN_ICON_STYLE = { color: 'rgb(62, 134, 53)' };

type IconConfig = {
  Icon: React.ComponentType<{ className?: string; style?: React.CSSProperties }>;
  label: string;
  green?: boolean;
};

const iconMap: Record<McpCardIconType, IconConfig> = {
  [McpCardIconType.VERIFIED_SOURCE]: {
    Icon: ShieldAltIcon,
    label: 'Verified source',
    green: true,
  },
  [McpCardIconType.SECURE_ENDPOINT]: {
    Icon: ShieldAltIcon,
    label: 'Secure endpoint',
    green: true,
  },
  [McpCardIconType.SAST]: {
    Icon: CheckCircleIcon,
    label: 'SAST',
    green: true,
  },
  [McpCardIconType.LOCAL_TO_CLUSTER]: {
    Icon: ServerIcon,
    label: 'Local to cluster',
    green: false,
  },
  [McpCardIconType.READ_ONLY_TOOLS]: {
    Icon: ToolsIcon,
    label: 'Read only tools',
    green: true,
  },
  [McpCardIconType.RED_HAT_PARTNER]: {
    Icon: CheckCircleIcon,
    label: 'Red Hat partner',
    green: true,
  },
  [McpCardIconType.REMOTE]: {
    Icon: ServerIcon,
    label: 'Remote',
    green: false,
  },
};

export const getMcpCardIconConfig = (type: McpCardIconType): (typeof iconMap)[McpCardIconType] =>
  iconMap[type];

export const getMcpCardIconConfigByLabel = (label: string): IconConfig | null => {
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
