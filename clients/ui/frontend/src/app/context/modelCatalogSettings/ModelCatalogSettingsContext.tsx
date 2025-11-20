import * as React from 'react';

export type ModelCatalogSettingsContextType = Record<string, never>;

type ModelCatalogSettingsContextProviderProps = {
  children: React.ReactNode;
};

export const ModelCatalogSettingsContext = React.createContext<ModelCatalogSettingsContextType>({});

export const ModelCatalogSettingsContextProvider: React.FC<
  ModelCatalogSettingsContextProviderProps
> = ({ children }) => {
  const contextValue = React.useMemo(() => ({}), []);

  return (
    <ModelCatalogSettingsContext.Provider value={contextValue}>
      {children}
    </ModelCatalogSettingsContext.Provider>
  );
};
