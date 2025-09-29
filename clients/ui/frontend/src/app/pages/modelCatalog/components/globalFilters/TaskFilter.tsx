import * as React from 'react';
import {
  ModelCatalogFilterResponseType,
  ModelCatalogFilterStatesByKey,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_TASK_NAME_MAPPING,
} from '~/concepts/modelCatalog/const';

const filterKey = ModelCatalogFilterKeys.TASK;

type TaskFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const TaskFilter: React.FC<TaskFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const task = filters[filterKey];

  React.useEffect(() => {
    const state: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
      task.values.forEach((key) => {
        state[key] = false;
      });
    if (task && !(filterKey in filterData)) {
      const state: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
      task.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData(filterKey, state);
    } else if (task) {
      setFilterData(filterKey, {
        ...filterData,
        [filterKey]: task.values.map((key) => ({
          key,
          value: false,
        })),
      });
    }
  }, [task, filterData, setFilterData]);

  if (!task) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.TASK>
      title="Task"
      filterToNameMapping={MODEL_CATALOG_TASK_NAME_MAPPING}
      filters={task}
      data={filterData[filterKey]}
      setData={(state) => setFilterData(filterKey, state)}
    />
  );
};

export default TaskFilter;
