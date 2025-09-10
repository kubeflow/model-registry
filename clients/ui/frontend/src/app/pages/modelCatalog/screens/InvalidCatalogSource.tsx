import { EmptyStateErrorMessage } from 'mod-arch-shared';
import React from 'react';

type InvalidCatalogSourceProps = {
  title?: string;
  sourceId?: string;
};

const InvalidCatalogSource: React.FC<InvalidCatalogSourceProps> = ({ title, sourceId }) => (
  <EmptyStateErrorMessage
    title={title || 'Source not found'}
    bodyText={`${sourceId ? `Catalog source ${sourceId}` : `The catalog source`} was not found`}
  />
);

export default InvalidCatalogSource;
