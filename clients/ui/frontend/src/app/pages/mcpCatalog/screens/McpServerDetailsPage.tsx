import * as React from 'react';
import { useParams } from 'react-router';
import { Link } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  Content,
  ContentVariants,
  Flex,
  FlexItem,
  Stack,
  StackItem,
  Spinner,
  Bullseye,
} from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import { mcpCatalogUrl } from '~/app/routes/mcpCatalog/mcpCatalog';
import { McpServerDetailsParams } from '~/app/routes/mcpCatalog/mcpServerDetails';
import { useMcpServerById } from '~/app/hooks/mcpCatalog/useMcpServerById';
import { useMcpCatalog } from '~/app/context/mcpCatalog/McpCatalogContext';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import McpServerDetailsView from './McpServerDetailsView';

const MCP_SERVER_ICON_URL =
  'https://catalog.redhat.com/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Funvalidated-model-logo.4b98e427.svg&w=96&q=75';

const McpServerDetailsPage: React.FC = () => {
  const { serverId } = useParams<keyof McpServerDetailsParams>();
  const decodedServerId = serverId ? decodeURIComponent(serverId) : '';

  const { apiState } = useMcpCatalog();
  const [server, loaded, loadError] = useMcpServerById(apiState, decodedServerId);

  // Check if server was found (id would be empty if not found/default)
  const serverNotFound = loaded && !loadError && !server.id;

  // Show loading state
  if (!loaded) {
    return (
      <ApplicationsPage
        breadcrumb={
          <Breadcrumb>
            <BreadcrumbItem>
              <Link to={mcpCatalogUrl()}>MCP catalog</Link>
            </BreadcrumbItem>
            <BreadcrumbItem isActive>Loading...</BreadcrumbItem>
          </Breadcrumb>
        }
        title="Loading..."
        empty={false}
        loaded={false}
        provideChildrenPadding
      >
        <Bullseye>
          <Spinner size="xl" />
        </Bullseye>
      </ApplicationsPage>
    );
  }

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        breadcrumb={
          <Breadcrumb>
            <BreadcrumbItem>
              <Link to={mcpCatalogUrl()}>MCP catalog</Link>
            </BreadcrumbItem>
            <BreadcrumbItem isActive>{server.name || 'Details'}</BreadcrumbItem>
          </Breadcrumb>
        }
        title={
          server.id ? (
            <Flex
              spaceItems={{ default: 'spaceItemsMd' }}
              alignItems={{ default: 'alignItemsCenter' }}
            >
              <img
                src={MCP_SERVER_ICON_URL}
                alt="MCP Server"
                style={{ height: '56px', width: '56px' }}
              />
              <Stack>
                <StackItem>
                  <Flex
                    spaceItems={{ default: 'spaceItemsSm' }}
                    alignItems={{ default: 'alignItemsCenter' }}
                  >
                    <FlexItem>{server.name}</FlexItem>
                  </Flex>
                </StackItem>
                <StackItem>
                  <Content component={ContentVariants.small}>
                    Provided by {server.provider || 'Unknown'}
                  </Content>
                </StackItem>
              </Stack>
            </Flex>
          ) : null
        }
        empty={serverNotFound}
        emptyStatePage={
          serverNotFound ? (
            <div>
              Server not found. Return to <Link to={mcpCatalogUrl()}>MCP catalog</Link>
            </div>
          ) : undefined
        }
        loadError={loadError}
        loaded={loaded}
        errorMessage="Unable to load MCP server"
        provideChildrenPadding
      >
        {server.id && <McpServerDetailsView server={server} />}
      </ApplicationsPage>
    </>
  );
};

export default McpServerDetailsPage;
