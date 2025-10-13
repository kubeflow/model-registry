import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type ModelCatalogLabelsProps = {
  tasks?: string[];
  license?: string;
  provider?: string;
  labels?: string[];
  numLabels: number;
};

const ModelCatalogLabels: React.FC<ModelCatalogLabelsProps> = ({
  tasks = [],
  license,
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
    {provider && <Label isCompact>{provider}</Label>}
    {labels.map((label) => (
      <Label data-testid="model-catalog-label" key={label} variant="outline">
        {label}
      </Label>
    ))}
    {license && <Label color="purple">{license}</Label>}
    {provider && <Label>{provider}</Label>}
  </LabelGroup>
);

export default ModelCatalogLabels;
