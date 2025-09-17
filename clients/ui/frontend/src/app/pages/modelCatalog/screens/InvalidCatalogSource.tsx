import { EmptyStateErrorMessage } from 'mod-arch-shared';
import React from 'react';
import { modelCatalogUrl } from '~/app/routes/modelCatalog/catalogModel';
import ModelCatalogSourceSelectorNavigator from './ModelCatalogSourceSelectorNavigator';

type InvalidCatalogSourceProps = {
  title?: string;
  sourceId?: string;
};

const InvalidCatalogSource: React.FC<InvalidCatalogSourceProps> = ({ title, sourceId }) => (
  <EmptyStateErrorMessage
    title={title || 'Source not found'}
    bodyText={`${sourceId ? `Catalog source ${sourceId}` : `The catalog source`} was not found`}
  >
    <ModelCatalogSourceSelectorNavigator
      getRedirectPath={(id: string) => modelCatalogUrl(id)}
      isPrimary
    />
  </EmptyStateErrorMessage>
);

export default InvalidCatalogSource;
