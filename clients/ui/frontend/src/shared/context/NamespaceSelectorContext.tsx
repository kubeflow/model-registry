import * as React from 'react';
import useNamespaces from '~/shared/hooks/useNamespaces';
import { Namespace } from '~/shared/types';

export type NamespaceSelectorContextType = {
  namespacesLoaded: boolean;
  namespacesLoadError?: Error;
  namespaces: Namespace[];
  preferredNamespace: Namespace | undefined;
  updatePreferredNamespace: (namespace: Namespace | undefined) => void;
};

type NamespaceSelectorContextProviderProps = {
  children: React.ReactNode;
};

export const NamespaceSelectorContext = React.createContext<NamespaceSelectorContextType>({
  namespacesLoaded: false,
  namespacesLoadError: undefined,
  namespaces: [],
  preferredNamespace: undefined,
  updatePreferredNamespace: () => undefined,
});

export const NamespaceSelectorContextProvider: React.FC<NamespaceSelectorContextProviderProps> = ({
  children,
  ...props
}) => (
  <EnabledNamespaceSelectorContextProvider {...props}>
    {children}
  </EnabledNamespaceSelectorContextProvider>
);

const EnabledNamespaceSelectorContextProvider: React.FC<NamespaceSelectorContextProviderProps> = ({
  children,
}) => {
  const [namespaces, isLoaded, error] = useNamespaces();
  const [preferredNamespace, setPreferredNamespace] =
    React.useState<NamespaceSelectorContextType['preferredNamespace']>(undefined);

  const firstNamespace = namespaces.length > 0 ? namespaces[0] : null;

  const contextValue = React.useMemo(
    () => ({
      namespacesLoaded: isLoaded,
      namespacesLoadError: error,
      namespaces,
      preferredNamespace: preferredNamespace ?? firstNamespace ?? undefined,
      updatePreferredNamespace: setPreferredNamespace,
    }),
    [isLoaded, error, namespaces, preferredNamespace, firstNamespace],
  );

  return (
    <NamespaceSelectorContext.Provider value={contextValue}>
      {children}
    </NamespaceSelectorContext.Provider>
  );
};
