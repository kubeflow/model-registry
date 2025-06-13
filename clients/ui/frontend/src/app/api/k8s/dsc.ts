import { k8sListResource, k8sPatchResource } from '@openshift/dynamic-plugin-sdk-utils';
import { DataScienceClusterKind, DeploymentMode } from '~/app/k8sTypes';
import { DataScienceClusterModel } from '~/app/api/models';

export const listDataScienceClusters = (): Promise<DataScienceClusterKind[]> =>
  k8sListResource<DataScienceClusterKind>({
    model: DataScienceClusterModel,
  }).then((dataScienceClusters) => dataScienceClusters.items);

export const patchDefaultDeploymentMode = (
  deploymentMode: DeploymentMode,
  dscName: string,
): Promise<DataScienceClusterKind> =>
  k8sPatchResource<DataScienceClusterKind>({
    model: DataScienceClusterModel,
    queryOptions: { name: dscName },
    patches: [
      {
        op: 'replace',
        path: '/spec/components/kserve/defaultDeploymentMode',
        value: deploymentMode,
      },
    ],
  }); 