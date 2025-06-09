import { Theme, DeploymentMode, PlatformMode, asEnumMember } from 'mod-arch-shared';

export const isIntegrated = (): boolean => DEPLOYMENT_MODE === DeploymentMode.Integrated;
export const isPlatformKubeflow = (): boolean => PLATFORM_MODE === PlatformMode.Kubeflow;

const STYLE_THEME = asEnumMember(process.env.STYLE_THEME, Theme) || Theme.Default;
const PLATFORM_MODE = asEnumMember(process.env.PLATFORM_MODE, PlatformMode) || PlatformMode.Default;
const DEPLOYMENT_MODE =
  asEnumMember(process.env.DEPLOYMENT_MODE, DeploymentMode) || DeploymentMode.Integrated;
const DEV_MODE = process.env.APP_ENV === 'development';
const POLL_INTERVAL = process.env.POLL_INTERVAL ? parseInt(process.env.POLL_INTERVAL) : 30000;
const AUTH_HEADER = process.env.AUTH_HEADER || 'kubeflow-userid';
const KUBEFLOW_USERNAME = process.env.KUBEFLOW_USERNAME || 'user@example.com';
const IMAGE_DIR = process.env.IMAGE_DIR || 'images';
const LOGO_LIGHT = process.env.LOGO || 'logo-light-theme.svg';
const URL_PREFIX = '/model-registry';
const BFF_API_VERSION = 'v1';

export {
  STYLE_THEME,
  POLL_INTERVAL,
  DEV_MODE,
  AUTH_HEADER,
  KUBEFLOW_USERNAME,
  IMAGE_DIR,
  LOGO_LIGHT,
  URL_PREFIX,
  PLATFORM_MODE,
  DEPLOYMENT_MODE,
  BFF_API_VERSION,
};

export const FindAdministratorOptions = [
  'The person who gave you your username, or who helped you to log in for the first time',
  'Someone in your IT department or help desk',
  'A project manager or developer',
];
