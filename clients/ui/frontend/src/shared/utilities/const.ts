export enum Theme {
  Default = 'default-theme',
  MUI = 'mui-theme',
  // Future themes can be added here
}

export enum DeploymentMode {
  Standalone = 'standalone',
  Integrated = 'integrated',
}

export enum PlatformMode {
  Default = 'default',
  Kubeflow = 'kubeflow',
}

export const isMUITheme = (): boolean => STYLE_THEME === Theme.MUI;
export const isStandalone = (): boolean => DEPLOYMENT_MODE === DeploymentMode.Standalone;
export const isIntegrated = (): boolean => DEPLOYMENT_MODE === DeploymentMode.Integrated;

export const isPlatformKubeflow = (): boolean => PLATFORM_MODE === PlatformMode.Kubeflow;
export const isPlatformDefault = (): boolean => PLATFORM_MODE === PlatformMode.Default;

const STYLE_THEME = process.env.STYLE_THEME || Theme.Default;
const PLATFORM_MODE = process.env.PLATFORM_MODE || PlatformMode.Default;
const DEV_MODE = process.env.APP_ENV === 'development';
const MOCK_AUTH = process.env.MOCK_AUTH === 'true';
const DEPLOYMENT_MODE = process.env.DEPLOYMENT_MODE || DeploymentMode.Integrated;
const POLL_INTERVAL = process.env.POLL_INTERVAL ? parseInt(process.env.POLL_INTERVAL) : 30000;
const AUTH_HEADER = process.env.AUTH_HEADER || 'kubeflow-userid';
const USERNAME = process.env.USERNAME || 'user@example.com';
const IMAGE_DIR = process.env.IMAGE_DIR || 'images';
const LOGO_LIGHT = process.env.LOGO || 'logo-light-theme.svg';
const URL_PREFIX = DEPLOYMENT_MODE === DeploymentMode.Integrated ? '/model-registry' : '';

export {
  POLL_INTERVAL,
  DEV_MODE,
  AUTH_HEADER,
  USERNAME,
  IMAGE_DIR,
  LOGO_LIGHT,
  MOCK_AUTH,
  URL_PREFIX,
  PLATFORM_MODE,
};

export const FindAdministratorOptions = [
  'The person who gave you your username, or who helped you to log in for the first time',
  'Someone in your IT department or help desk',
  'A project manager or developer',
];
