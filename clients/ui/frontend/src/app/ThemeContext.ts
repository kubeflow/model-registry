import * as React from 'react';

type ThemeContextProps = {
  isMUITheme: boolean;
};

export const ThemeContext = React.createContext({
  isMUITheme: true,
});

export const useThemeContext = (): ThemeContextProps => React.useContext(ThemeContext);
