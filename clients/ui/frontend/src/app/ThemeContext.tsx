import * as React from 'react';
import { createTheme } from '@mui/material';
import { ThemeProvider as MUIThemeProvider } from '@mui/material/styles';
import { isMUITheme, Theme } from '~/shared/utilities/const';

type ThemeProviderProps = {
  children: React.ReactNode;
};

type ThemeContextProps = {
  isMUITheme: boolean;
};

export const ThemeContext = React.createContext({
  isMUITheme: false,
});

export const useThemeContext = (): ThemeContextProps => React.useContext(ThemeContext);

const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const themeValue = React.useMemo(() => ({ isMUITheme: isMUITheme() }), []);
  const createMUITheme = React.useMemo(() => createTheme({ cssVariables: true }), []);

  React.useEffect(() => {
    // Apply the theme based on the value of STYLE_THEME
    if (isMUITheme()) {
      document.documentElement.classList.add(Theme.MUI);
    } else {
      document.documentElement.classList.remove(Theme.MUI);
    }
  }, []);

  return (
    <ThemeContext.Provider value={themeValue}>
      {isMUITheme() ? (
        <MUIThemeProvider theme={createMUITheme}>{children}</MUIThemeProvider>
      ) : (
        children
      )}
    </ThemeContext.Provider>
  );
};

export default ThemeProvider;
