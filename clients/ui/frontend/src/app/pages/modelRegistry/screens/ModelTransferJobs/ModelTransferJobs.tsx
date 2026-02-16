import React from 'react';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';
import { Link, useParams } from 'react-router-dom';
import { ApplicationsPage } from 'mod-arch-shared';
import { ModelTransferJob } from '~/app/types';
import useModelTransferJobs from '~/app/hooks/useModelTransferJobs';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import {
  modelRegistryUrl,
  modelTransferJobsUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import ModelRegistrySelectorNavigator from '~/app/pages/modelRegistry/screens/ModelRegistrySelectorNavigator';
import DeleteModelTransferJobModal from './DeleteModelTransferJobModal';
import ModelTransferJobsListView from './ModelTransferJobsListView';

type ModelTransferJobsProps = Omit<
  React.ComponentProps<typeof ApplicationsPage>,
  'breadcrumb' | 'title' | 'description' | 'loadError' | 'loaded' | 'provideChildrenPadding'
>;

const ModelTransferJobs: React.FC<ModelTransferJobsProps> = ({ ...pageProps }) => {
  const { modelRegistry } = useParams<{ modelRegistry: string }>();
  const [jobs, jobsLoaded, jobsLoadError, refetchJobs] = useModelTransferJobs();
  const { api, apiAvailable } = useModelRegistryAPI();
  const [jobToDelete, setJobToDelete] = React.useState<ModelTransferJob | null>(null);

  const onDeleteTransferJob = React.useCallback(
    async (job: ModelTransferJob) => {
      if (!apiAvailable) {
        throw new Error('API not available');
      }
      await api.deleteModelTransferJob({}, job.id);
      setJobToDelete(null);
      await refetchJobs();
    },
    [api, apiAvailable, refetchJobs],
  );

  const onRequestDelete = React.useCallback((job: ModelTransferJob) => {
    setJobToDelete(job);
  }, []);

  const onCloseDeleteModal = React.useCallback(() => setJobToDelete(null), []);

  return (
    <ApplicationsPage
      {...pageProps}
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem
            render={() => <Link to={modelRegistryUrl(modelRegistry)}>Model registry</Link>}
          />
          <BreadcrumbItem data-testid="breadcrumb-transfer-jobs" isActive>
            Model transfer jobs
          </BreadcrumbItem>
        </Breadcrumb>
      }
      title="Model transfer jobs"
      description="Monitor the status of model transfer jobs. Model transfer jobs are created when you choose to store model artifacts at registration time. When a job is complete, the registered model version appears in the specified model registry."
      loadError={jobsLoadError}
      loaded={jobsLoaded}
      provideChildrenPadding
      headerContent={
        <ModelRegistrySelectorNavigator
          getRedirectPath={(modelRegistryName) => modelTransferJobsUrl(modelRegistryName)}
        />
      }
    >
      <ModelTransferJobsListView jobs={jobs.items} onRequestDelete={onRequestDelete} />
      {jobToDelete && (
        <DeleteModelTransferJobModal
          job={jobToDelete}
          onClose={onCloseDeleteModal}
          onDelete={onDeleteTransferJob}
        />
      )}
    </ApplicationsPage>
  );
};

export default ModelTransferJobs;
