export enum Theme {
  Default = 'default-theme',
  MUI = 'mui-theme',
  // Future themes can be added here
}

export enum DeploymentMode {
  Standalone = 'standalone',
  Integrated = 'integrated',
}

export const isMUITheme = (): boolean => STYLE_THEME === Theme.MUI;
export const isStandalone = (): boolean => DEPLOYMENT_MODE === DeploymentMode.Standalone;
export const isIntegrated = (): boolean => DEPLOYMENT_MODE === DeploymentMode.Integrated;

const STYLE_THEME = process.env.STYLE_THEME || Theme.MUI;
const DEV_MODE = process.env.APP_ENV === 'development';
const MOCK_AUTH = process.env.MOCK_AUTH === 'true';
const DEPLOYMENT_MODE = process.env.DEPLOYMENT_MODE || DeploymentMode.Integrated;
const POLL_INTERVAL = process.env.POLL_INTERVAL ? parseInt(process.env.POLL_INTERVAL) : 30000;
const AUTH_HEADER = process.env.AUTH_HEADER || 'kubeflow-userid';
const USERNAME = process.env.USERNAME || 'user@example.com';
const IMAGE_DIR = process.env.IMAGE_DIR || 'images';
const LOGO_LIGHT = process.env.LOGO || 'logo-light-theme.svg';

export { POLL_INTERVAL, DEV_MODE, AUTH_HEADER, USERNAME, IMAGE_DIR, LOGO_LIGHT, MOCK_AUTH };
