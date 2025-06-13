import { ConfigSecretItem } from '~/app/k8sTypes';
import { SecureDBInfo, SecureDBRType } from '~/app/pages/modelRegistrySettings/const';

export const isClusterWideCABundleEnabled = (configMaps: ConfigSecretItem[]): boolean =>
  configMaps.some(
    (configMap) =>
      configMap.name === 'odh-trusted-ca-bundle' &&
      configMap.keys.some((key: string) => key === 'odh-ca-bundle.crt'),
  );

export const isOpenshiftCAbundleEnabled = (configMaps: ConfigSecretItem[]): boolean =>
  configMaps.some(
    (configMap) =>
      configMap.name === 'openshift-service-ca.crt' &&
      configMap.keys.some((key: string) => key === 'service-ca.crt'),
  );

export const findSecureDBType = (resourceName: string, key: string): SecureDBRType => {
  if (resourceName === 'odh-trusted-ca-bundle' && key === 'odh-ca-bundle.crt') {
    return SecureDBRType.CLUSTER_WIDE;
  }
  if (resourceName === 'openshift-service-ca.crt' && key === 'service-ca.crt') {
    return SecureDBRType.OPENSHIFT;
  }
  return SecureDBRType.EXISTING;
};

export const findConfigMap = (
  secureDBInfo: SecureDBInfo,
): { name: string; key: string } | null => {
  if (secureDBInfo.type === SecureDBRType.CLUSTER_WIDE) {
    return { name: 'odh-trusted-ca-bundle', key: 'odh-ca-bundle.crt' };
  }
  if (secureDBInfo.type === SecureDBRType.OPENSHIFT) {
    return { name: 'openshift-service-ca.crt', key: 'service-ca.crt' };
  }
  return { name: secureDBInfo.resourceName, key: secureDBInfo.key };
}; 