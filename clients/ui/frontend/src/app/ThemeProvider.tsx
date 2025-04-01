import * as React from 'react';
import { createTheme } from '@mui/material';
import { ThemeProvider as MUIThemeProvider } from '@mui/material/styles';
import { isMUITheme, Theme } from '~/shared/utilities/const';
import { ThemeContext } from './ThemeContext';

type ThemeProviderProps = {
  children: React.ReactNode;
};

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
