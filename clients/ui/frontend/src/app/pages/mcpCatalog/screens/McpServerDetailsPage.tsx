import React from 'react';
import { useParams } from 'react-router';
import { Link, useLocation } from 'react-router-dom';
import {
  ActionList,
  ActionListGroup,
  Breadcrumb,
  BreadcrumbItem,
  Button,
  Content,
  ContentVariants,
  Flex,
  FlexItem,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { ApplicationsIcon } from '@patternfly/react-icons';
import { ApplicationsPage } from 'mod-arch-shared';
import { useMcpServer } from '~/app/hooks/mcpServerCatalog/useMcpServer';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';
import { mcpCatalogUrl } from '~/app/routes/mcpCatalog/mcpCatalog';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import McpServerDetailsView from './McpServerDetailsView';

const McpServerDetailsPage: React.FC = () => {
  const { serverId = '' } = useParams<{ serverId: string }>();
  const location = useLocation();
  const [apiServer, apiServerLoaded, apiServerLoadError] = useMcpServer(serverId);
  const { mcpServers, mcpServersLoaded } = React.useContext(McpCatalogContext);

  const { server, serverLoaded } = React.useMemo(() => {
    if (apiServerLoaded && apiServer && !apiServerLoadError) {
      return { server: apiServer, serverLoaded: true };
    }

    const contextMatch = mcpServers.items.find((s) => String(s.id) === serverId);
    if (contextMatch) {
      return { server: contextMatch, serverLoaded: true };
    }

    const mockMatch = mockMcpServers.find((s) => String(s.id) === serverId);
    if (mockMatch) {
      return { server: mockMatch, serverLoaded: true };
    }

    if (apiServerLoaded || mcpServersLoaded) {
      return { server: undefined, serverLoaded: true };
    }

    return { server: undefined, serverLoaded: false };
  }, [
    apiServer,
    apiServerLoaded,
    apiServerLoadError,
    mcpServers.items,
    mcpServersLoaded,
    serverId,
  ]);

  const catalogLink = React.useMemo(() => {
    const searchParams = new URLSearchParams(location.search);
    const base = mcpCatalogUrl();
    return searchParams.toString() ? `${base}?${searchParams.toString()}` : base;
  }, [location.search]);

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        breadcrumb={
          <Breadcrumb>
            <BreadcrumbItem>
              <Link to={catalogLink}>MCP Catalog</Link>
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
            <div>
              Details not found. Return to <Link to={mcpCatalogUrl()}>MCP Catalog</Link>
            </div>
          ) : undefined
        }
        loaded={serverLoaded}
        provideChildrenPadding
        headerAction={
          serverLoaded &&
          server && (
            <ActionList>
              <ActionListGroup>
                <Button
                  variant="primary"
                  data-testid="deploy-mcp-server-button"
                  onClick={() => {
                    // Stub handler for future story
                  }}
                >
                  Deploy MCP Server
                </Button>
              </ActionListGroup>
            </ActionList>
          )
        }
      >
        {server && <McpServerDetailsView server={server} />}
      </ApplicationsPage>
    </>
  );
};

export default McpServerDetailsPage;
