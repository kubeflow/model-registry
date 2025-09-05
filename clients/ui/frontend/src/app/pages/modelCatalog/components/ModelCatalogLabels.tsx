import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type ModelCatalogLabelsProps = {
  tags?: string[];
  labels?: string[];
  tasks?: string[];
  license?: string;
  provider?: string;
};

const ModelCatalogLabels: React.FC<ModelCatalogLabelsProps> = ({
  tags,
  labels,
  tasks = [],
  license,
  provider,
}) => (
  <LabelGroup numLabels={5} isCompact>
    {labels && (
      <Label color="blue" isCompact>
        {labels}
      </Label>
    )}
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
    {tags?.map((tag) => (
      <Label key={tag} color="grey" isCompact>
        {tag}
      </Label>
    ))}
  </LabelGroup>
);

export default ModelCatalogLabels;
