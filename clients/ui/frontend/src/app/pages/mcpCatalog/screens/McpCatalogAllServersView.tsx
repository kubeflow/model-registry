import React from 'react';
import { Stack } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import {
  filterEnabledCatalogSources,
  hasSourcesWithoutLabels,
  orderLabelsByPriority,
  getUniqueSourceLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { SourceLabel } from '~/app/modelCatalogTypes';
import McpCatalogCategorySection from './McpCatalogCategorySection';

type McpCatalogAllServersViewProps = {
  searchTerm: string;
};

const McpCatalogAllServersView: React.FC<McpCatalogAllServersViewProps> = ({ searchTerm }) => {
  const { catalogSources, catalogLabels, setSelectedSourceLabel } =
    React.useContext(McpCatalogContext);

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);
    return orderLabelsByPriority(uniqueLabels, catalogLabels);
  }, [catalogSources, catalogLabels]);

  const hasSourcesWithoutLabelsValue = React.useMemo(
    () => hasSourcesWithoutLabels(catalogSources),
    [catalogSources],
  );

  const handleShowMoreCategory = React.useCallback(
    (categoryLabel: string) => {
      setSelectedSourceLabel(categoryLabel);
    },
    [setSelectedSourceLabel],
  );

  return (
    <Stack hasGutter>
      {sourceLabels.map((label) => (
        <McpCatalogCategorySection
          key={label}
          label={label}
          searchTerm={searchTerm}
          pageSize={4}
          onShowMore={handleShowMoreCategory}
        />
      ))}
      {hasSourcesWithoutLabelsValue && (
        <McpCatalogCategorySection
          key="other-servers"
          label={SourceLabel.other}
          searchTerm={searchTerm}
          pageSize={4}
          onShowMore={handleShowMoreCategory}
        />
      )}
    </Stack>
  );
};

export default McpCatalogAllServersView;
