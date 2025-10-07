import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogStringFilterKey,
  MODEL_CATALOG_TASK_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptions, ModelCatalogStringFilterOptions } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogStringFilterKey.TASK;

type TaskFilterProps = {
  filters?: Extract<CatalogFilterOptions, Partial<ModelCatalogStringFilterOptions>>;
};

const TaskFilter: React.FC<TaskFilterProps> = ({ filters }) => {
  const task = filters?.[filterKey];
  if (!task) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogStringFilterKey.TASK>
      title="Task"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_TASK_NAME_MAPPING}
      filters={task}
    />
  );
};

export default TaskFilter;
