import * as React from 'react';
import { Alert, Bullseye } from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import { Outlet } from 'react-router-dom';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';

const MCP_CATALOG_TITLE = 'MCP Catalog';
const MCP_CATALOG_DESCRIPTION =
  'Discover and manage MCP servers and tools available for your organization.';

const McpCatalogCoreLoader: React.FC = () => {
  const { catalogSourcesLoaded, catalogSourcesLoadError } = React.useContext(McpCatalogContext);

  if (catalogSourcesLoadError) {
    return (
      <ApplicationsPage
        title={MCP_CATALOG_TITLE}
        description={MCP_CATALOG_DESCRIPTION}
        headerContent={null}
        empty
        emptyStatePage={
          <Bullseye>
            <Alert title="MCP catalog source load error" variant="danger" isInline>
              {catalogSourcesLoadError.message}
            </Alert>
          </Bullseye>
        }
        loaded
      />
    );
  }

  if (!catalogSourcesLoaded) {
    return (
      <ApplicationsPage
        title={MCP_CATALOG_TITLE}
        description={MCP_CATALOG_DESCRIPTION}
        headerContent={null}
        empty
        emptyStatePage={<Bullseye>Loading catalog sources...</Bullseye>}
        loaded={false}
      />
    );
  }

  return <Outlet />;
};

export default McpCatalogCoreLoader;
