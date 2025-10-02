import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKey,
  MODEL_CATALOG_TASK_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, GlobalFilterTypes } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKey.TASK;

type TaskFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<GlobalFilterTypes>>;
};

const TaskFilter: React.FC<TaskFilterProps> = ({ filters }) => {
  const task = filters?.[filterKey];
  if (!task) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKey.TASK>
      title="Task"
      filterKey={filterKey}
      filterToNameMapping={MODEL_CATALOG_TASK_NAME_MAPPING}
      filters={task}
    />
  );
};

export default TaskFilter;
