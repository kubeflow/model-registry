import * as React from 'react';
import useNamespaces from '~/shared/hooks/useNamespaces';
import { Namespace } from '~/shared/types';
import { isIntegrated } from '~/shared/utilities/const';

export type NamespaceSelectorContextType = {
  namespacesLoaded: boolean;
  namespacesLoadError?: Error;
  namespaces: Namespace[];
  preferredNamespace: Namespace | undefined;
  updatePreferredNamespace: (namespace: Namespace | undefined) => void;
  initializationError?: Error;
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
  initializationError: undefined,
});

export const NamespaceSelectorContextProvider: React.FC<NamespaceSelectorContextProviderProps> = ({
  children,
  ...props
}) => (
  <EnabledNamespaceSelectorContextProvider {...props}>
    {children}
  </EnabledNamespaceSelectorContextProvider>
);

declare global {
  interface Window {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    centraldashboard: any;
  }
}

const EnabledNamespaceSelectorContextProvider: React.FC<NamespaceSelectorContextProviderProps> = ({
  children,
}) => {
  const [unsortedNamespaces, isLoaded, error] = useNamespaces();
  const namespaces = React.useMemo(
    () => unsortedNamespaces.toSorted((a, b) => a.name.localeCompare(b.name)),
    [unsortedNamespaces],
  );
  const [preferredNamespace, setPreferredNamespace] =
    React.useState<NamespaceSelectorContextType['preferredNamespace']>(undefined);
  const [initializationError, setInitializationError] = React.useState<Error>();

  const firstNamespace = namespaces.length > 0 ? namespaces[0] : null;

  React.useEffect(() => {
    if (isIntegrated()) {
      // Initialize the central dashboard client
      try {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        window.centraldashboard.CentralDashboardEventHandler.init((cdeh: any) => {
          // eslint-disable-next-line no-param-reassign
          cdeh.onNamespaceSelected = (newNamespace: string) => {
            setPreferredNamespace({ name: newNamespace });
          };
        });
      } catch (err) {
        /* eslint-disable no-console */
        console.error('Failed to initialize central dashboard client', err);
        if (err instanceof Error) {
          setInitializationError(err);
        }
      }
    }
  }, []);

  const contextValue = React.useMemo(
    () => ({
      namespacesLoaded: isLoaded,
      namespacesLoadError: error,
      namespaces,
      preferredNamespace: preferredNamespace ?? firstNamespace ?? undefined,
      updatePreferredNamespace: setPreferredNamespace,
      initializationError,
    }),
    [isLoaded, error, namespaces, preferredNamespace, firstNamespace, initializationError],
  );

  return (
    <NamespaceSelectorContext.Provider value={contextValue}>
      {children}
    </NamespaceSelectorContext.Provider>
  );
};
