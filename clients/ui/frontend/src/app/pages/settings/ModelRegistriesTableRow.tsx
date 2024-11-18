import React from 'react';
import { Td, Tr } from '@patternfly/react-table';
import { ModelRegistry } from '~/app/types';

type ModelRegistriesTableRowProps = {
  modelRegistry: ModelRegistry;
};

const ModelRegistriesTableRow: React.FC<ModelRegistriesTableRowProps> = ({ modelRegistry: mr }) => (
  <>
    <Tr>
      <Td dataLabel="Model registry name">
        <strong>{mr.displayName || mr.name}</strong>
        {mr.description && <p>{mr.description}</p>}
      </Td>
    </Tr>
  </>
);

// TODO: [Model Registry RBAC] Get rest of columns once we manage permissions

export default ModelRegistriesTableRow;
