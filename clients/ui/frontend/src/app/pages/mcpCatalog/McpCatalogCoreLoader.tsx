import * as React from 'react';
import { Alert, Bullseye } from '@patternfly/react-core';
import {
  ApplicationsPage,
  KubeflowDocs,
  ProjectObjectType,
  TitleWithIcon,
  typedEmptyImage,
  WhosMyAdministrator,
} from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import { Outlet } from 'react-router-dom';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';

const MCP_CATALOG_TITLE = 'MCP Catalog';
const MCP_CATALOG_DESCRIPTION =
  'Discover and manage MCP servers and tools available for your organization.';

const McpCatalogCoreLoader: React.FC = () => {
  const { catalogSourcesLoaded, catalogSourcesLoadError, sourceLabels } =
    React.useContext(McpCatalogContext);
  const { isMUITheme } = useThemeContext();

  if (catalogSourcesLoadError) {
    return (
      <ApplicationsPage
        title={
          <TitleWithIcon title={MCP_CATALOG_TITLE} objectType={ProjectObjectType.modelCatalog} />
        }
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
        title={
          <TitleWithIcon title={MCP_CATALOG_TITLE} objectType={ProjectObjectType.modelCatalog} />
        }
        description={MCP_CATALOG_DESCRIPTION}
        headerContent={null}
        empty
        emptyStatePage={<Bullseye>Loading catalog sources...</Bullseye>}
        loaded={false}
      />
    );
  }

  if (sourceLabels.length === 0) {
    return (
      <ApplicationsPage
        title={
          <TitleWithIcon title={MCP_CATALOG_TITLE} objectType={ProjectObjectType.modelCatalog} />
        }
        description={MCP_CATALOG_DESCRIPTION}
        empty
        emptyStatePage={
          <EmptyModelCatalogState
            testid="empty-mcp-catalog-state"
            title="MCP catalog configuration required"
            description={
              isMUITheme
                ? 'To discover MCP servers, follow the instructions in the docs below.'
                : 'There are no MCP sources to display. Request that your administrator configure model sources for the catalog.'
            }
            headerIcon={() => (
              <img src={typedEmptyImage(ProjectObjectType.modelRegistrySettings)} alt="" />
            )}
            primaryAction={isMUITheme ? <KubeflowDocs /> : <WhosMyAdministrator />}
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

export default McpCatalogCoreLoader;
