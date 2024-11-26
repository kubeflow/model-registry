import React from 'react';
import { NamespaceSelectorContext } from '~/shared/context/NamespaceSelectorContext';
import { isStandalone } from '~/shared/utilities/const';
import { useDeepCompareMemoize } from '~/shared/utilities/useDeepCompareMemoize';

const useQueryParamNamespaces = (): Record<string, unknown> => {
  const { preferredNamespace: namespaceSelector } = React.useContext(NamespaceSelectorContext);
  // TODO: Readd GetNamespaceQueryParam once it's working
  const namespace = isStandalone() ? namespaceSelector?.name : 'kubeflow';

  return useDeepCompareMemoize({ namespace });
};

export default useQueryParamNamespaces;
