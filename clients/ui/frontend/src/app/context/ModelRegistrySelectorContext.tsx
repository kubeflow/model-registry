import * as React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import { ModelRegistry } from '~/app/types';
import useModelRegistries from '~/app/hooks/useModelRegistries';

export type ModelRegistrySelectorContextType = {
  modelRegistriesLoaded: boolean;
  modelRegistriesLoadError?: Error;
  modelRegistries: ModelRegistry[];
  preferredModelRegistry: ModelRegistry | undefined;
  updatePreferredModelRegistry: (modelRegistry: ModelRegistry | undefined) => void;
  //refreshRulesReview: () => void; TODO: [Midstream] Reimplement this
};

type ModelRegistrySelectorContextProviderProps = {
  children: React.ReactNode;
};

export const ModelRegistrySelectorContext = React.createContext<ModelRegistrySelectorContextType>({
  modelRegistriesLoaded: false,
  modelRegistriesLoadError: undefined,
  modelRegistries: [],
  preferredModelRegistry: undefined,
  updatePreferredModelRegistry: () => undefined,
  //refreshRulesReview: () => undefined,
});

export const ModelRegistrySelectorContextProvider: React.FC<
  ModelRegistrySelectorContextProviderProps
> = ({ children, ...props }) => (
  <EnabledModelRegistrySelectorContextProvider {...props}>
    {children}
  </EnabledModelRegistrySelectorContextProvider>
);

const EnabledModelRegistrySelectorContextProvider: React.FC<
  ModelRegistrySelectorContextProviderProps
> = ({ children }) => {
  // TODO: [Midstream] Add area check for enablement

  const queryParams = useQueryParamNamespaces();

  const [modelRegistries, isLoaded, error] = useModelRegistries(queryParams);
  const [preferredModelRegistry, setPreferredModelRegistry] =
    React.useState<ModelRegistrySelectorContextType['preferredModelRegistry']>(undefined);

  const firstModelRegistry = modelRegistries.length > 0 ? modelRegistries[0] : null;

  const contextValue = React.useMemo(
    () => ({
      modelRegistriesLoaded: isLoaded,
      modelRegistriesLoadError: error,
      modelRegistries,
      preferredModelRegistry: preferredModelRegistry ?? firstModelRegistry ?? undefined,
      updatePreferredModelRegistry: setPreferredModelRegistry,
      // refreshRulesReview,
    }),
    [isLoaded, error, modelRegistries, preferredModelRegistry, firstModelRegistry],
  );

  return (
    <ModelRegistrySelectorContext.Provider value={contextValue}>
      {children}
    </ModelRegistrySelectorContext.Provider>
  );
};
