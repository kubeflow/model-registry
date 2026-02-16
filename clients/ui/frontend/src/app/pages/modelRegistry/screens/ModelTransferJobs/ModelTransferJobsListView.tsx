import * as React from 'react';
import { ToolbarGroup } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { ModelTransferJob } from '~/app/types';
import EmptyModelRegistryState from '~/app/pages/modelRegistry/screens/components/EmptyModelRegistryState';
import FilterToolbar from '~/app/shared/components/FilterToolbar';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import ModelTransferJobsTable from './ModelTransferJobsTable';

enum ModelTransferJobsFilterOptions {
  jobName = 'jobName',
  modelName = 'modelName',
  versionName = 'versionName',
  namespace = 'namespace',
  author = 'author',
  status = 'status',
}

const modelTransferJobsFilterOptions = {
  [ModelTransferJobsFilterOptions.jobName]: 'Job name',
  [ModelTransferJobsFilterOptions.modelName]: 'Model name',
  [ModelTransferJobsFilterOptions.versionName]: 'Version name',
  [ModelTransferJobsFilterOptions.namespace]: 'Namespace',
  [ModelTransferJobsFilterOptions.author]: 'Author',
  [ModelTransferJobsFilterOptions.status]: 'Status',
};

type ModelTransferJobsFilterDataType = Record<ModelTransferJobsFilterOptions, string | undefined>;

const initialFilterData: ModelTransferJobsFilterDataType = {
  [ModelTransferJobsFilterOptions.jobName]: undefined,
  [ModelTransferJobsFilterOptions.modelName]: undefined,
  [ModelTransferJobsFilterOptions.versionName]: undefined,
  [ModelTransferJobsFilterOptions.namespace]: undefined,
  [ModelTransferJobsFilterOptions.author]: undefined,
  [ModelTransferJobsFilterOptions.status]: undefined,
};

type ModelTransferJobsListViewProps = {
  jobs: ModelTransferJob[];
  onRequestDelete?: (job: ModelTransferJob) => void;
};

const ModelTransferJobsListView: React.FC<ModelTransferJobsListViewProps> = ({
  jobs,
  onRequestDelete,
}) => {
  const [filterData, setFilterData] =
    React.useState<ModelTransferJobsFilterDataType>(initialFilterData);

  const onFilterUpdate = React.useCallback(
    (key: string, value: string | { label: string; value: string } | undefined) =>
      setFilterData((prevValues) => ({ ...prevValues, [key]: value })),
    [setFilterData],
  );

  const onClearFilters = React.useCallback(() => setFilterData(initialFilterData), [setFilterData]);

  // Filter jobs based on all filter criteria
  const filteredJobs = React.useMemo(() => {
    const jobNameFilter = filterData[ModelTransferJobsFilterOptions.jobName]?.toLowerCase();
    const modelNameFilter = filterData[ModelTransferJobsFilterOptions.modelName]?.toLowerCase();
    const versionNameFilter = filterData[ModelTransferJobsFilterOptions.versionName]?.toLowerCase();
    const namespaceFilter = filterData[ModelTransferJobsFilterOptions.namespace]?.toLowerCase();
    const authorFilter = filterData[ModelTransferJobsFilterOptions.author]?.toLowerCase();
    const statusFilter = filterData[ModelTransferJobsFilterOptions.status]?.toLowerCase();

    return jobs.filter((job) => {
      if (jobNameFilter && !job.name.toLowerCase().includes(jobNameFilter)) {
        return false;
      }
      if (modelNameFilter && !job.registeredModelName?.toLowerCase().includes(modelNameFilter)) {
        return false;
      }
      if (versionNameFilter && !job.modelVersionName?.toLowerCase().includes(versionNameFilter)) {
        return false;
      }
      if (namespaceFilter && !job.namespace?.toLowerCase().includes(namespaceFilter)) {
        return false;
      }
      if (authorFilter && !job.author?.toLowerCase().includes(authorFilter)) {
        return false;
      }
      if (statusFilter && !job.status.toLowerCase().includes(statusFilter)) {
        return false;
      }
      return true;
    });
  }, [jobs, filterData]);

  if (jobs.length === 0) {
    return (
      <EmptyModelRegistryState
        testid="empty-model-transfer-jobs"
        title="No model transfer jobs"
        headerIcon={SearchIcon}
        description="Model transfer jobs are created when you choose to store model artifacts during model registration."
      />
    );
  }

  const toggleGroupItems = (
    <ToolbarGroup variant="filter-group">
      <FilterToolbar
        filterOptions={modelTransferJobsFilterOptions}
        filterOptionRenders={{
          [ModelTransferJobsFilterOptions.jobName]: ({ onChange, ...props }) => (
            <ThemeAwareSearchInput
              {...props}
              fieldLabel="Filter by job name"
              placeholder="Filter by job name"
              className="toolbar-fieldset-wrapper"
              style={{ minWidth: '270px' }}
              onChange={(value) => onChange(value)}
            />
          ),
          [ModelTransferJobsFilterOptions.modelName]: ({ onChange, ...props }) => (
            <ThemeAwareSearchInput
              {...props}
              fieldLabel="Filter by model name"
              placeholder="Filter by model name"
              className="toolbar-fieldset-wrapper"
              style={{ minWidth: '270px' }}
              onChange={(value) => onChange(value)}
            />
          ),
          [ModelTransferJobsFilterOptions.versionName]: ({ onChange, ...props }) => (
            <ThemeAwareSearchInput
              {...props}
              fieldLabel="Filter by version name"
              placeholder="Filter by version name"
              className="toolbar-fieldset-wrapper"
              style={{ minWidth: '270px' }}
              onChange={(value) => onChange(value)}
            />
          ),
          [ModelTransferJobsFilterOptions.namespace]: ({ onChange, ...props }) => (
            <ThemeAwareSearchInput
              {...props}
              fieldLabel="Filter by namespace"
              placeholder="Filter by namespace"
              className="toolbar-fieldset-wrapper"
              style={{ minWidth: '270px' }}
              onChange={(value) => onChange(value)}
            />
          ),
          [ModelTransferJobsFilterOptions.author]: ({ onChange, ...props }) => (
            <ThemeAwareSearchInput
              {...props}
              fieldLabel="Filter by author"
              placeholder="Filter by author"
              className="toolbar-fieldset-wrapper"
              style={{ minWidth: '270px' }}
              onChange={(value) => onChange(value)}
            />
          ),
          [ModelTransferJobsFilterOptions.status]: ({ onChange, ...props }) => (
            <ThemeAwareSearchInput
              {...props}
              fieldLabel="Filter by status"
              placeholder="Filter by status"
              className="toolbar-fieldset-wrapper"
              style={{ minWidth: '270px' }}
              onChange={(value) => onChange(value)}
            />
          ),
        }}
        filterData={filterData}
        onFilterUpdate={onFilterUpdate}
      />
    </ToolbarGroup>
  );

  return (
    <ModelTransferJobsTable
      jobs={filteredJobs}
      clearFilters={onClearFilters}
      toolbarContent={toggleGroupItems}
      onRequestDelete={onRequestDelete}
    />
  );
};

export default ModelTransferJobsListView;
