import React from 'react';
import { Stack } from '@patternfly/react-core';
import { useMcpCatalog } from '~/app/context/mcpCatalog/McpCatalogContext';
import {
  filterEnabledMcpSources,
  getUniqueMcpSourceLabels,
  hasMcpSourcesWithoutLabels,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import { McpCategoryName, McpSourceLabel } from '~/app/pages/mcpCatalog/types';
import McpCatalogCategorySection from '~/app/pages/mcpCatalog/components/McpCatalogCategorySection';

/**
 * All Servers view showing servers grouped by source label categories.
 * Uses server-side filtering via context - no props needed for filters.
 */
const McpCatalogAllServersView: React.FC = () => {
  const {
    mcpSources,
    mcpServers,
    mcpServersLoaded,
    mcpServersLoadError,
    updateSelectedSourceLabel,
  } = useMcpCatalog();

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledMcpSources(mcpSources);
    return getUniqueMcpSourceLabels(enabledSources);
  }, [mcpSources]);

  const hasSourcesWithoutLabelsValue = React.useMemo(
    () => hasMcpSourcesWithoutLabels(mcpSources),
    [mcpSources],
  );

  const handleShowMoreCategory = React.useCallback(
    (categoryLabel: string) => {
      updateSelectedSourceLabel(categoryLabel);
    },
    [updateSelectedSourceLabel],
  );

  const servers = mcpServers?.items || [];

  return (
    <Stack hasGutter>
      {sourceLabels.map((label) => (
        <McpCatalogCategorySection
          key={label}
          label={label}
          pageSize={4}
          servers={servers}
          sources={mcpSources}
          loaded={mcpServersLoaded}
          loadError={mcpServersLoadError}
          onShowMore={handleShowMoreCategory}
        />
      ))}
      {hasSourcesWithoutLabelsValue && (
        <McpCatalogCategorySection
          key={McpCategoryName.communityAndCustomServers}
          label={McpSourceLabel.other}
          pageSize={4}
          servers={servers}
          sources={mcpSources}
          loaded={mcpServersLoaded}
          loadError={mcpServersLoadError}
          onShowMore={handleShowMoreCategory}
          displayName={McpCategoryName.communityAndCustomServers}
        />
      )}
    </Stack>
  );
};

export default McpCatalogAllServersView;
