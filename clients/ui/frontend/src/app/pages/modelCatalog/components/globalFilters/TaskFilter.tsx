import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_TASK_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, ModelCatalogFilterTypesByKey } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKeys.TASK;

type TaskFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<ModelCatalogFilterTypesByKey>>;
};

const TaskFilter: React.FC<TaskFilterProps> = ({ filters }) => {
  const task = filters?.[filterKey];
  if (!task) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.TASK>
      title="Task"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_TASK_NAME_MAPPING}
      filters={task}
    />
  );
};

export default TaskFilter;
