import * as React from 'react';

type AreaContextState = {
  dscStatus: unknown;
};

export const AreaContext = React.createContext<AreaContextState>({
  dscStatus: {
    components: {
      modelregistry: {
        registriesNamespace: 'opendatahub',
      },
    },
  },
});

const AreaContextProvider: React.FC<React.PropsWithChildren> = ({ children }) => (
  <AreaContext.Provider
    value={React.useMemo(
      () => ({
        dscStatus: { components: { modelregistry: { registriesNamespace: 'opendatahub' } } },
      }),
      [],
    )}
  >
    {children}
  </AreaContext.Provider>
);
export default AreaContextProvider;
