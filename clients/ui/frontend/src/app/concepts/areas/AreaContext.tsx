import * as React from 'react';

type AreaContextState = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  // placeholder for future context values
};

export const AreaContext = React.createContext<AreaContextState>({
  // TODO: Add default values here
});

const AreaContextProvider: React.FC<React.PropsWithChildren> = ({ children }) => (
  <AreaContext.Provider value={React.useMemo(() => ({}), [])}>{children}</AreaContext.Provider>
);
export default AreaContextProvider;
