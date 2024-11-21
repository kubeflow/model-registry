// TODO: [Env Handling] Fetch the .env variable here.
const POLL_INTERVAL = 30000;

export enum Theme {
  Default = 'default-theme',
  MUI = 'mui-theme',
  // Future themes can be added here
}

export const isMUITheme = (): boolean => STYLE_THEME === Theme.MUI;

export const STYLE_THEME = process.env.STYLE_THEME || Theme.MUI;

export const USER_ACCESS_TOKEN = 'x-forwarded-access-token';

export { POLL_INTERVAL };
