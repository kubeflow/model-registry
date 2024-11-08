import * as React from 'react';
import { UserSettings, ConfigSettings } from '~/shared/types';

type AppContextProps = {
  config: ConfigSettings;
  user: UserSettings;
};

// eslint-disable-next-line @typescript-eslint/consistent-type-assertions
export const AppContext = React.createContext({} as AppContextProps);

export const useAppContext = (): AppContextProps => React.useContext(AppContext);
