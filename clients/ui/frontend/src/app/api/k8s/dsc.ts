import { DataScienceClusterKind, DeploymentMode } from '~/app/k8sTypes';

export const listDataScienceClusters = (): Promise<DataScienceClusterKind[]> =>
  Promise.resolve([
    {
        apiVersion: 'datasciencecluster.opendatahub.io/v1',
        kind: 'DataScienceCluster',
        metadata: {
            name: 'default',
            namespace: 'opendatahub'
        },
        spec: {
            components: {
                modelregistry: {
                    registriesNamespace: 'opendatahub'
                }
            }
        },
        status: {
            conditions: [],
            installedComponents: {},
            phase: 'Ready',
            release: {
                name: '2.5.0',
                version: '2.5.0'
            }
        }
    }
  ]);

export const patchDefaultDeploymentMode = (
  deploymentMode: DeploymentMode,
  dscName: string,
): Promise<DataScienceClusterKind> =>
  Promise.resolve({
    apiVersion: 'datasciencecluster.opendatahub.io/v1',
    kind: 'DataScienceCluster',
    metadata: {
        name: 'default',
        namespace: 'opendatahub'
    },
    spec: {
        components: {
            modelregistry: {
                registriesNamespace: 'opendatahub'
            }
        }
    },
    status: {
        conditions: [],
        installedComponents: {},
        phase: 'Ready',
        release: {
            name: '2.5.0',
            version: '2.5.0'
        }
    }
  }); 