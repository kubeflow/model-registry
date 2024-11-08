import * as React from 'react';
import EmptyStateErrorMessage from '~/shared/components/EmptyStateErrorMessage';
import { modelRegistryUrl } from './routeUtils';
import ModelRegistrySelectorNavigator from './ModelRegistrySelectorNavigator';

type InvalidModelRegistryProps = {
  title?: string;
  modelRegistry?: string;
};

const InvalidModelRegistry: React.FC<InvalidModelRegistryProps> = ({ title, modelRegistry }) => (
  <EmptyStateErrorMessage
    title={title || 'Model Registry not found'}
    bodyText={`${
      modelRegistry ? `Model Registry ${modelRegistry}` : 'The Model Registry'
    } was not found.`}
  >
    <ModelRegistrySelectorNavigator
      getRedirectPath={(modelRegistryName) => modelRegistryUrl(modelRegistryName)}
      primary
    />
  </EmptyStateErrorMessage>
);

export default InvalidModelRegistry;
