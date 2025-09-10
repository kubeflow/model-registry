import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type ModelCatalogLabelsProps = {
  tasks?: string[];
  license?: string;
  provider?: string;
};

const ModelCatalogLabels: React.FC<ModelCatalogLabelsProps> = ({
  tasks = [],
  license,
  provider,
}) => (
  <LabelGroup numLabels={5} isCompact>
    {tasks.map((task) => (
      <Label data-testid="model-catalog-label" key={task} variant="outline">
        {task}
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
