import * as React from 'react';
import { Label, LabelGroup } from '@patternfly/react-core';

type ModelCatalogLabelsProps = {
  tags?: string[];
  framework?: string;
  task?: string;
  license?: string;
};

const ModelCatalogLabels: React.FC<ModelCatalogLabelsProps> = ({
  tags,
  framework,
  task,
  license,
}) => (
  <LabelGroup numLabels={5} isCompact>
    {framework && (
      <Label color="blue" isCompact>
        {framework}
      </Label>
    )}
    {task && (
      <Label color="green" isCompact>
        {task}
      </Label>
    )}
    {license && (
      <Label color="purple" isCompact>
        {license}
      </Label>
    )}
    {tags?.map((tag) => (
      <Label key={tag} color="grey" isCompact>
        {tag}
      </Label>
    ))}
  </LabelGroup>
);

export default ModelCatalogLabels;
