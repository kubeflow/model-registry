import * as React from 'react';
import { useQueryParamNamespaces, useBrowserStorage } from 'mod-arch-core';
import { ModelRegistry } from '~/app/types';
import useModelRegistries from '~/app/hooks/useModelRegistries';

export const MODEL_REGISTRY_SELECTED_STORAGE_KEY = 'kubeflow.dashboard.model.registry.selected';
export const MODEL_REGISTRY_FAVORITE_STORAGE_KEY = 'kubeflow.dashboard.model.registry.favorite';

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

  // Persist selected registry name in sessionStorage so it survives in-session navigation
  const [preferredModelRegistryName, setPreferredModelRegistryName] = useBrowserStorage<string>(
    MODEL_REGISTRY_SELECTED_STORAGE_KEY,
    '',
    true,
    true, // use sessionStorage
  );

  // Read favorites from localStorage (same key used by the selector component)
  const [favorites] = useBrowserStorage<string[]>(MODEL_REGISTRY_FAVORITE_STORAGE_KEY, []);

  // Resolve the preferred registry object from the persisted name
  const preferredModelRegistry = preferredModelRegistryName
    ? (modelRegistries.find((mr) => mr.name === preferredModelRegistryName) ?? undefined)
    : undefined;

  // Fallback: first favorite registry, then first available registry
  const firstFavoriteRegistry =
    favorites.length > 0
      ? (modelRegistries.find((mr) => favorites.includes(mr.name)) ?? null)
      : null;
  const firstModelRegistry = modelRegistries.length > 0 ? modelRegistries[0] : null;

  const updatePreferredModelRegistry = React.useCallback(
    (modelRegistry: ModelRegistry | undefined) => {
      setPreferredModelRegistryName(modelRegistry?.name ?? '');
    },
    [setPreferredModelRegistryName],
  );

  const contextValue = React.useMemo(
    () => ({
      modelRegistriesLoaded: isLoaded,
      modelRegistriesLoadError: error,
      modelRegistries,
      preferredModelRegistry:
        preferredModelRegistry ?? firstFavoriteRegistry ?? firstModelRegistry ?? undefined,
      updatePreferredModelRegistry,
      // refreshRulesReview,
    }),
    [
      isLoaded,
      error,
      modelRegistries,
      preferredModelRegistry,
      firstFavoriteRegistry,
      firstModelRegistry,
      updatePreferredModelRegistry,
    ],
  );

  return (
    <ModelRegistrySelectorContext.Provider value={contextValue}>
      {children}
    </ModelRegistrySelectorContext.Provider>
  );
};
