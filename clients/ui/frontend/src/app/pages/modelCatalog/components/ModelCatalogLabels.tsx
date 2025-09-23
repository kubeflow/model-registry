import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type ModelCatalogLabelsProps = {
  tasks?: string[];
  license?: string;
  provider?: string;
  labels?: string[];
};

const ModelCatalogLabels: React.FC<ModelCatalogLabelsProps> = ({
  tasks = [],
  license,
  provider,
  labels = [],
}) => (
  <LabelGroup numLabels={5} isCompact>
    {tasks.map((task) => (
      <Label data-testid="model-catalog-label" key={task} variant="outline">
        {task}
      </Label>
    ))}
    {labels.map((label) => (
      <Label data-testid="model-catalog-label" key={label} variant="outline">
        {label}
      </Label>
    ))}
    {license && (
      <Label color="purple" isCompact>
        {license}
      </Label>
    )}
    {provider && <Label isCompact>{provider}</Label>}
  </LabelGroup>
);

export default ModelCatalogLabels;
