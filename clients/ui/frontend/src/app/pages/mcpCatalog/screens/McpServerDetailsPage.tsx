import React from 'react';
import { useParams } from 'react-router';
import { Link } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  Button,
  Content,
  ContentVariants,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  Flex,
  FlexItem,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { ApplicationsIcon, SearchIcon } from '@patternfly/react-icons';
import { ApplicationsPage } from 'mod-arch-shared';
import { useMcpServer } from '~/app/hooks/mcpServerCatalog/useMcpServer';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';
import { mcpCatalogUrl } from '~/app/routes/mcpCatalog/mcpCatalog';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import McpServerDetailsView from './McpServerDetailsView';

const McpServerDetailsPage: React.FC = () => {
  const { serverId = '' } = useParams<{ serverId: string }>();
  const [apiServer, apiServerLoaded, apiServerLoadError] = useMcpServer(serverId);
  const { mcpServers, mcpServersLoaded } = React.useContext(McpCatalogContext);

  const { server, serverLoaded, serverLoadError } = React.useMemo(() => {
    if (apiServerLoaded && apiServer && !apiServerLoadError) {
      return { server: apiServer, serverLoaded: true, serverLoadError: undefined };
    }

    const contextMatch = mcpServers.items.find((s) => String(s.id) === serverId);
    if (contextMatch) {
      return { server: contextMatch, serverLoaded: true, serverLoadError: undefined };
    }

    const mockMatch = mockMcpServers.find((s) => String(s.id) === serverId);
    if (mockMatch) {
      return { server: mockMatch, serverLoaded: true, serverLoadError: undefined };
    }

    if (apiServerLoaded || mcpServersLoaded) {
      // Context loads all servers; if it loaded and still no match, it's genuinely not found.
      // Only propagate API errors when the context hasn't loaded yet (real failures like 500s).
      const isNotFound = mcpServersLoaded;
      return {
        server: undefined,
        serverLoaded: true,
        serverLoadError: isNotFound ? undefined : apiServerLoadError,
      };
    }

    return { server: undefined, serverLoaded: false, serverLoadError: undefined };
  }, [
    apiServer,
    apiServerLoaded,
    apiServerLoadError,
    mcpServers.items,
    mcpServersLoaded,
    serverId,
  ]);

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        breadcrumb={
          <Breadcrumb>
            <BreadcrumbItem>
              <Link to={mcpCatalogUrl()}>MCP Catalog</Link>
            </BreadcrumbItem>
            <BreadcrumbItem isActive data-testid="breadcrumb-server-name">
              {server?.name || 'Details'}
            </BreadcrumbItem>
          </Breadcrumb>
        }
        title={
          server ? (
            <Flex
              spaceItems={{ default: 'spaceItemsMd' }}
              alignItems={{ default: 'alignItemsCenter' }}
            >
              {server.logo ? (
                <img
                  src={server.logo}
                  alt="server logo"
                  style={{ height: '56px', width: '56px' }}
                />
              ) : (
                <ApplicationsIcon
                  style={{ fontSize: '56px' }}
                  data-testid="mcp-server-default-icon"
                />
              )}
              <Stack>
                <StackItem>
                  <FlexItem>{server.name}</FlexItem>
                </StackItem>
                {server.provider && (
                  <StackItem>
                    <Content component={ContentVariants.small}>Provider: {server.provider}</Content>
                  </StackItem>
                )}
              </Stack>
            </Flex>
          ) : null
        }
        empty={!server}
        emptyStatePage={
          !server ? (
            <EmptyState
              icon={SearchIcon}
              titleText="MCP server not found"
              data-testid="mcp-server-not-found"
            >
              <EmptyStateBody>The requested MCP server could not be found.</EmptyStateBody>
              <EmptyStateFooter>
                <Button
                  variant="primary"
                  component={(props) => <Link {...props} to={mcpCatalogUrl()} />}
                >
                  Return to MCP Catalog
                </Button>
              </EmptyStateFooter>
            </EmptyState>
          ) : undefined
        }
        loadError={serverLoadError}
        loaded={serverLoaded}
        errorMessage="Unable to load MCP server details"
        provideChildrenPadding
      >
        {server && <McpServerDetailsView server={server} />}
      </ApplicationsPage>
    </>
  );
};

export default McpServerDetailsPage;
