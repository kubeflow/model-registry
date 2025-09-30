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
  const currentState = filterData[filterKey];

  React.useEffect(() => {
    if (!task) {
      return;
    }

    const filterKeys = task.values;
    const hasMatchingKeys =
      currentState !== undefined &&
      filterKeys.length === Object.keys(currentState).length &&
      filterKeys.every((key) => key in currentState);

    if (hasMatchingKeys) {
      return;
    }

    const nextState: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
    filterKeys.forEach((key) => {
      nextState[key] = currentState?.[key] ?? false;
    });

    setFilterData(filterKey, nextState);
  }, [task, currentState, setFilterData]);

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
