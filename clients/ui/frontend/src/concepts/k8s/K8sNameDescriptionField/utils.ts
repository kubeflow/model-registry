import {
  K8sNameDescriptionFieldData,
  K8sNameDescriptionFieldUpdateFunctionInternal,
  UseK8sNameDescriptionDataConfiguration,
} from './types';

const MAX_K8S_NAME_LENGTH = 253;

/**
 * Translates a name to a k8s-safe value.
 * @see https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
 */
export const translateDisplayNameForK8s = (name = '', safePrefix = ''): string => {
  const translatedName = name
    .trim()
    .toLowerCase()
    .replace(/\s/g, '-')
    .replace(/[^A-Za-z0-9-]/g, '');

  if (safePrefix) {
    return `${safePrefix}${translatedName}`;
  }

  return translatedName;
};

export const checkValidK8sName = (
  value: string,
): { valid: boolean; invalidCharacters: boolean } => {
  if (!value) {
    return { valid: true, invalidCharacters: false };
  }

  // Kubernetes name must consist of lower case alphanumeric characters or '-'
  // and must start and end with an alphanumeric character
  const valid = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/.test(value);
  const invalidCharacters = !/^[a-z0-9-]*$/.test(value);

  return { valid, invalidCharacters };
};

export const setupDefaults = (
  configuration: UseK8sNameDescriptionDataConfiguration,
): K8sNameDescriptionFieldData => {
  const { initialData, editableK8sName } = configuration;

  const name = initialData?.name ?? '';
  const description = initialData?.description ?? '';
  const k8sName = initialData?.k8sName ?? '';

  return {
    name,
    description,
    k8sName: {
      value: k8sName,
      state: {
        immutable: !!(k8sName && !editableK8sName),
        invalidCharacters: false,
        invalidLength: k8sName.length > MAX_K8S_NAME_LENGTH,
        maxLength: MAX_K8S_NAME_LENGTH,
        touched: !!k8sName,
      },
    },
  };
};

export const handleUpdateLogic =
  (currentData: K8sNameDescriptionFieldData): K8sNameDescriptionFieldUpdateFunctionInternal =>
  (key, value) => {
    if (key === 'k8sName') {
      const validation = checkValidK8sName(value);
      return {
        ...currentData,
        k8sName: {
          value,
          state: {
            ...currentData.k8sName.state,
            invalidCharacters: !validation.valid && validation.invalidCharacters,
            invalidLength: value.length > MAX_K8S_NAME_LENGTH,
            touched: true,
          },
        },
      };
    }

    if (key === 'name') {
      const k8sValue = currentData.k8sName.state.touched
        ? currentData.k8sName.value
        : translateDisplayNameForK8s(value);

      const validation = checkValidK8sName(k8sValue);

      return {
        ...currentData,
        name: value,
        k8sName: {
          value: k8sValue,
          state: {
            ...currentData.k8sName.state,
            invalidCharacters: !validation.valid && validation.invalidCharacters,
            invalidLength: k8sValue.length > MAX_K8S_NAME_LENGTH,
          },
        },
      };
    }

    return {
      ...currentData,
      [key]: value,
    };
  };
