import * as React from 'react';
import { StackItem } from '@patternfly/react-core';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import type { McpCatalogFilterOptions } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

const filterKey = 'securityVerification';

type SecurityVerificationFilterProps = {
  filters?: McpCatalogFilterOptions;
};

const SecurityVerificationFilter: React.FC<SecurityVerificationFilterProps> = ({ filters }) => {
  const value = filters?.[filterKey];
  if (!value) {
    return null;
  }
  return (
    <StackItem>
      <McpCatalogStringFilter
        title="Security & Verification"
        filterKey={filterKey}
        filters={value}
      />
    </StackItem>
  );
};

export default SecurityVerificationFilter;
