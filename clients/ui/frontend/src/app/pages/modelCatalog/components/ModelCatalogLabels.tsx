import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type ModelCatalogLabelsProps = {
  tasks?: string[];
  provider?: string;
  labels?: string[];
  numLabels: number;
};

const ModelCatalogLabels: React.FC<ModelCatalogLabelsProps> = ({
  tasks = [],
  provider,
  labels = [],
  numLabels,
}) => (
  <LabelGroup numLabels={numLabels} isCompact>
    {tasks.map((task) => (
      <Label data-testid="model-catalog-label" key={task} variant="outline">
        {task}
      </Label>
    ))}
    {provider && (
      <Label isCompact variant="outline">
        {provider}
      </Label>
    )}
    {labels.map((label) => (
      <Label data-testid="model-catalog-label" key={label} variant="outline">
        {label}
      </Label>
    ))}
  </LabelGroup>
);

export default ModelCatalogLabels;
