import * as React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import useModelRegistryAPIState, {
  ModelRegistryAPIState,
} from '~/app/hooks/useModelRegistryAPIState';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';

export type ModelRegistryContextType = {
  apiState: ModelRegistryAPIState;
  refreshAPIState: () => void;
};

type ModelRegistryContextProviderProps = {
  children: React.ReactNode;
  modelRegistryName: string;
};

export const ModelRegistryContext = React.createContext<ModelRegistryContextType>({
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  apiState: { apiAvailable: false, api: null as unknown as ModelRegistryAPIState['api'] },
  refreshAPIState: () => undefined,
});

export const ModelRegistryContextProvider: React.FC<ModelRegistryContextProviderProps> = ({
  children,
  modelRegistryName,
}) => {
  const hostPath = modelRegistryName
    ? `${URL_PREFIX}/api/${BFF_API_VERSION}/model_registry/${modelRegistryName}`
    : null;

  const queryParams = useQueryParamNamespaces();

  const [apiState, refreshAPIState] = useModelRegistryAPIState(hostPath, queryParams);

  return (
    <ModelRegistryContext.Provider
      value={React.useMemo(
        () => ({
          apiState,
          refreshAPIState,
        }),
        [apiState, refreshAPIState],
      )}
    >
      {children}
    </ModelRegistryContext.Provider>
  );
};
