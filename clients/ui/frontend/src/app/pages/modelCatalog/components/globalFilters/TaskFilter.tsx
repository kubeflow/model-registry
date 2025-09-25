import * as React from 'react';
import {
  ModelCatalogStringFilterStateType,
  ModelCatalogFilterResponseType,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

type TaskFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const TaskFilter: React.FC<TaskFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const { task } = filters;

  React.useEffect(() => {
    if (task && !('task' in filterData)) {
      const state: Record<string, boolean> = {};
      task.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData('task', state);
    }
  }, [task, filterData, setFilterData]);

  if (!task) {
    return null;
  }

  return (
    <ModelCatalogStringFilter
      title="Task"
      filterKey="task"
      filters={task}
      data={filterData}
      setData={(state: ModelCatalogStringFilterStateType) => setFilterData('task', state)}
    />
  );
};

export default TaskFilter;
