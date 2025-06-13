import * as React from 'react';
import { translateDisplayNameForK8s } from '~/app/concepts/k8s/utils';
import { K8sDSGResource } from '~/app/k8sTypes';

export type K8sNameDescriptionFieldType = {
  name: string;
  k8sName: {
    value: string;
    error: string;
  };
  description: string;
};

export const useK8sNameDescriptionFieldData = ({
  initialData,
  k8sNameIsEditable,
}: {
  initialData?: K8sDSGResource;
  k8sNameIsEditable?: boolean;
} = {}): {
  data: K8sNameDescriptionFieldType;
  onDataChange: (data: K8sNameDescriptionFieldType) => void;
} => {
  const [data, setData] = React.useState<K8sNameDescriptionFieldType>({
    name: initialData?.metadata.annotations?.['openshift.io/display-name'] || '',
    k8sName: {
      value: initialData?.metadata.name || '',
      error: '',
    },
    description: initialData?.metadata.annotations?.['openshift.io/description'] || '',
  });

  return {
    data,
    onDataChange: (newData) => {
      const k8sName = k8sNameIsEditable
        ? newData.k8sName.value
        : translateDisplayNameForK8s(newData.name);
      setData({ ...newData, k8sName: { value: k8sName, error: '' } });
    },
  };
}; 