export enum SecureDBRType {
    EXISTING = 'existing',
    NEW = 'new',
    CLUSTER_WIDE = 'cluster-wide',
    OPENSHIFT = 'openshift',
}

export enum ResourceType {
    Secret = 'Secret',
    ConfigMap = 'ConfigMap',
}

export type SecureDBInfo = {
    type: SecureDBRType;
    nameSpace: string;
    resourceName: string;
    certificate: string;
    key: string;
    isValid: boolean;
    resourceType?: ResourceType;
}; 