import React from 'react';
import { NamespaceSelectorContext } from '~/shared/context/NamespaceSelectorContext';
import { isStandalone } from '~/shared/utilities/const';
import { getNamespaceQueryParam } from '~/shared/api/apiUtils';
import { useDeepCompareMemoize } from '~/shared/utilities/useDeepCompareMemoize';

const useQueryParamNamespaces = (): Record<string, unknown> => {
  const { preferredNamespace: namespaceSelector } = React.useContext(NamespaceSelectorContext);
  const namespace = isStandalone() ? namespaceSelector?.name : getNamespaceQueryParam();

  return useDeepCompareMemoize({ namespace });
};

export default useQueryParamNamespaces;
