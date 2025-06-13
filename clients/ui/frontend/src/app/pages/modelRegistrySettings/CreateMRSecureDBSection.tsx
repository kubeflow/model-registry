import * as React from 'react';
import {
  FormGroup,
  TextField,
  Select,
  MenuItem,
  InputLabel,
  FormControl,
} from '@mui/material';
import { ResourceType, SecureDBRType } from './const';

type CreateMRSecureDBSectionProps = {
  secureDBInfo: SecureDBInfo;
  modelRegistryNamespace: string;
  k8sName: string;
  existingCertConfigMaps: { name: string; keys: string[] }[];
  existingCertSecrets: { name: string; keys: string[] }[];
  setSecureDBInfo: (info: SecureDBInfo) => void;
};

export const CreateMRSecureDBSection: React.FC<CreateMRSecureDBSectionProps> = ({
  secureDBInfo,
  modelRegistryNamespace,
  k8sName,
  existingCertConfigMaps,
  existingCertSecrets,
  setSecureDBInfo,
}) => (
  <>
    <FormControl fullWidth>
      <InputLabel id="db-type-label">Database Type</InputLabel>
      <Select
        labelId="db-type-label"
        value={secureDBInfo.type}
        label="Database Type"
        onChange={(e) => setSecureDBInfo({ ...secureDBInfo, type: e.target.value as SecureDBRType })}
      >
        <MenuItem value={SecureDBRType.EXISTING}>Existing</MenuItem>
        <MenuItem value={SecureDBRType.NEW}>New</MenuItem>
        <MenuItem value={SecureDBRType.CLUSTER_WIDE}>Cluster Wide</MenuItem>
        <MenuItem value={SecureDBRType.OPENSHIFT}>OpenShift</MenuItem>
      </Select>
    </FormControl>
    {secureDBInfo.type === SecureDBRType.EXISTING && (
      <FormControl fullWidth>
        <InputLabel id="resource-type-label">Resource Type</InputLabel>
        <Select
            labelId="resource-type-label"
            value={secureDBInfo.resourceType}
            label="Resource Type"
            onChange={(e) => setSecureDBInfo({ ...secureDBInfo, resourceType: e.target.value as ResourceType })}
        >
            <MenuItem value={ResourceType.Secret}>Secret</MenuItem>
            <MenuItem value={ResourceType.ConfigMap}>ConfigMap</MenuItem>
        </Select>
      </FormControl>
    )}
    {/* Add more fields based on the selected type */}
  </>
);

export type SecureDBInfo = any; 