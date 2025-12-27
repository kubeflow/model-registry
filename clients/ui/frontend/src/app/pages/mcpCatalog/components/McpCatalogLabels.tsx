import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type McpCatalogLabelsProps = {
  tags?: string[];
  provider?: string;
  numLabels?: number;
};

const McpCatalogLabels: React.FC<McpCatalogLabelsProps> = ({
  tags = [],
  provider,
  numLabels = 3,
}) => (
  <LabelGroup numLabels={numLabels} isCompact>
    {tags.map((tag) => (
      <Label data-testid="mcp-catalog-label" key={tag} variant="outline">
        {tag}
      </Label>
    ))}
    {provider && (
      <Label isCompact variant="outline">
        {provider}
      </Label>
    )}
  </LabelGroup>
);

export default McpCatalogLabels;
