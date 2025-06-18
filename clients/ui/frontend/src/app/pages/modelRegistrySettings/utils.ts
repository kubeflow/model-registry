// Generic, upstream-friendly utilities for Model Registry Secure DB handling
// No OpenShift/ODH specifics; all such values must be provided via config

import { ConfigSecretItem, ModelRegistryKind } from '~/app/k8sTypes';
import { RecursivePartial } from '~/typeHelpers';

export enum SecureDBRType {
  EXISTING = 'existing',
  CLUSTER_WIDE = 'clusterWide',
  OPENSHIFT = 'openshift',
  NEW = 'new',
}

export enum ResourceType {
  Secret = 'Secret',
  ConfigMap = 'ConfigMap',
}

export interface CABundleConfig {
  trustedBundleName: string;
  clusterWideKey: string;
  openshiftKey: string;
}

export interface SecureDBInfo {
  type: string;
  resourceType?: string;
  resourceName?: string;
  key?: string;
  // Add other fields as needed
}

export const findSecureDBType = (
  name: string,
  key: string,
  config: CABundleConfig,
): SecureDBRType => {
  if (name === config.trustedBundleName && key === config.clusterWideKey) {
    return SecureDBRType.CLUSTER_WIDE;
  }
  if (name === config.trustedBundleName && key === config.openshiftKey) {
    return SecureDBRType.OPENSHIFT;
  }
  return SecureDBRType.EXISTING;
};

export const findConfigMap = (
  secureDBInfo: SecureDBInfo,
  config: CABundleConfig,
): { name: string; key: string } => {
  if (secureDBInfo.type === SecureDBRType.CLUSTER_WIDE) {
    return { name: config.trustedBundleName, key: config.clusterWideKey };
  }
  if (secureDBInfo.type === SecureDBRType.OPENSHIFT) {
    return { name: config.trustedBundleName, key: config.openshiftKey };
  }
  return { name: secureDBInfo.resourceName ?? '', key: secureDBInfo.key ?? '' };
};

export const constructRequestBody = (
  data: RecursivePartial<ModelRegistryKind>,
  secureDBInfo: SecureDBInfo,
  addSecureDB: boolean,
  config: CABundleConfig,
): RecursivePartial<ModelRegistryKind> => {
  const mr = data;
  if (addSecureDB && secureDBInfo.resourceType === ResourceType.Secret && mr.spec?.mysql) {
    mr.spec.mysql.sslRootCertificateSecret = {
      name: secureDBInfo.resourceName ?? '',
      key: secureDBInfo.key ?? '',
    };
    mr.spec.mysql.sslRootCertificateConfigMap = null;
  } else if (addSecureDB && mr.spec?.mysql) {
    mr.spec.mysql.sslRootCertificateConfigMap = findConfigMap(secureDBInfo, config);
    mr.spec.mysql.sslRootCertificateSecret = null;
  } else if (!addSecureDB && mr.spec?.mysql) {
    mr.spec.mysql.sslRootCertificateConfigMap = null;
    mr.spec.mysql.sslRootCertificateSecret = null;
  }
  return mr;
};

export const isClusterWideCABundleEnabled = (
  existingCertConfigMaps: ConfigSecretItem[],
  config: CABundleConfig,
): boolean => {
  const clusterWideCABundle = existingCertConfigMaps.find(
    (configMap) =>
      configMap.name === config.trustedBundleName && configMap.keys.includes(config.clusterWideKey),
  );
  return !!clusterWideCABundle;
};

export const isOpenshiftCAbundleEnabled = (
  existingCertConfigMaps: ConfigSecretItem[],
  config: CABundleConfig,
): boolean => {
  const openshiftCAbundle = existingCertConfigMaps.find(
    (configMap) =>
      configMap.name === config.trustedBundleName && configMap.keys.includes(config.openshiftKey),
  );
  return !!openshiftCAbundle;
};
