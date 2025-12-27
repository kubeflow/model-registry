import { Alert, Bullseye } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import {
  ApplicationsPage,
  KubeflowDocs,
  ProjectObjectType,
  TitleWithIcon,
  typedEmptyImage,
  WhosMyAdministrator,
} from 'mod-arch-shared';
import * as React from 'react';
import { Outlet } from 'react-router-dom';
import {
  McpCatalogContext,
  McpCatalogContextProvider,
} from '~/app/context/mcpCatalog/McpCatalogContext';
import EmptyMcpCatalogState from './EmptyMcpCatalogState';

/**
 * McpCatalogCoreContent handles loading, error, and empty states for MCP Catalog sources.
 * Mirrors the pattern used in ModelCatalogCoreLoader for consistency.
 */
const McpCatalogCoreContent: React.FC = () => {
  const { mcpSources, mcpSourcesLoaded, mcpSourcesLoadError } = React.useContext(McpCatalogContext);

  const { isMUITheme } = useThemeContext();

  if (mcpSourcesLoadError) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="MCP Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover MCP servers that are available for your organization to deploy and use."
        headerContent={null}
        empty
        emptyStatePage={
          <Bullseye>
            <Alert title="MCP catalog source load error" variant="danger" isInline>
              {mcpSourcesLoadError.message}
            </Alert>
          </Bullseye>
        }
        loaded
      />
    );
  }

  if (!mcpSourcesLoaded) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="MCP Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover MCP servers that are available for your organization to deploy and use."
        headerContent={null}
        empty
        emptyStatePage={<Bullseye>Loading MCP catalog sources...</Bullseye>}
        loaded={false}
      />
    );
  }

  if (mcpSources?.items?.length === 0) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="MCP Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover MCP servers that are available for your organization to deploy and use."
        empty
        emptyStatePage={
          <EmptyMcpCatalogState
            testid="empty-mcp-catalog-state"
            title={isMUITheme ? 'Deploy an MCP catalog' : 'Request access to MCP catalog'}
            description={
              isMUITheme
                ? 'To deploy MCP catalog, follow the instructions in the docs below.'
                : 'To request MCP catalog, or to request permission to access MCP catalog, contact your administrator.'
            }
            headerIcon={() => (
              <img src={typedEmptyImage(ProjectObjectType.modelRegistrySettings)} alt="" />
            )}
            customAction={isMUITheme ? <KubeflowDocs /> : <WhosMyAdministrator />}
          />
        }
        headerContent={null}
        loaded
        provideChildrenPadding
      />
    );
  }

  return <Outlet />;
};

/**
 * McpCatalogCoreLoader wraps the MCP Catalog routes with the context provider
 * to provide API state and MCP server data to all child components.
 */
const McpCatalogCoreLoader: React.FC = () => (
  <McpCatalogContextProvider>
    <McpCatalogCoreContent />
  </McpCatalogContextProvider>
);

export default McpCatalogCoreLoader;
