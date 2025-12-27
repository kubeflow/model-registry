import * as React from 'react';
import { Divider, Stack, StackItem } from '@patternfly/react-core';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';

type McpCatalogFiltersProps = {
  allProviders: string[];
  allLicenses: string[];
  allTags: string[];
  allTransports: string[];
  allDeploymentModes: string[];
  selectedProviders: string[];
  selectedLicenses: string[];
  selectedTags: string[];
  selectedTransports: string[];
  selectedDeploymentModes: string[];
  onProviderChange: (provider: string, checked: boolean) => void;
  onLicenseChange: (license: string, checked: boolean) => void;
  onTagChange: (tag: string, checked: boolean) => void;
  onTransportChange: (transport: string, checked: boolean) => void;
  onDeploymentModeChange: (mode: string, checked: boolean) => void;
};

const McpCatalogFilters: React.FC<McpCatalogFiltersProps> = ({
  allProviders,
  allLicenses,
  allTags,
  allTransports,
  allDeploymentModes,
  selectedProviders,
  selectedLicenses,
  selectedTags,
  selectedTransports,
  selectedDeploymentModes,
  onProviderChange,
  onLicenseChange,
  onTagChange,
  onTransportChange,
  onDeploymentModeChange,
}) => (
  <Stack hasGutter>
    {allDeploymentModes.length > 0 && (
      <>
        <StackItem>
          <McpCatalogStringFilter
            title="Deployment Mode"
            values={allDeploymentModes}
            selectedValues={selectedDeploymentModes}
            onSelectionChange={onDeploymentModeChange}
          />
        </StackItem>
        <Divider />
      </>
    )}
    {allTransports.length > 0 && (
      <>
        <StackItem>
          <McpCatalogStringFilter
            title="Transport"
            values={allTransports}
            selectedValues={selectedTransports}
            onSelectionChange={onTransportChange}
          />
        </StackItem>
        <Divider />
      </>
    )}
    {allProviders.length > 0 && (
      <>
        <StackItem>
          <McpCatalogStringFilter
            title="Provider"
            values={allProviders}
            selectedValues={selectedProviders}
            onSelectionChange={onProviderChange}
          />
        </StackItem>
        <Divider />
      </>
    )}
    {allTags.length > 0 && (
      <>
        <StackItem>
          <McpCatalogStringFilter
            title="Tags"
            values={allTags}
            selectedValues={selectedTags}
            onSelectionChange={onTagChange}
          />
        </StackItem>
        <Divider />
      </>
    )}
    {allLicenses.length > 0 && (
      <StackItem>
        <McpCatalogStringFilter
          title="License"
          values={allLicenses}
          selectedValues={selectedLicenses}
          onSelectionChange={onLicenseChange}
        />
      </StackItem>
    )}
  </Stack>
);

export default McpCatalogFilters;
