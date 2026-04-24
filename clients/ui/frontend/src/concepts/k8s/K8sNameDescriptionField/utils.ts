import {
  K8sNameDescriptionFieldData,
  K8sNameDescriptionFieldUpdateFunctionInternal,
  UseK8sNameDescriptionDataConfiguration,
} from './types';

const MAX_K8S_NAME_LENGTH = 253;

/**
 * Deterministic [a-z0-9]+ suffix for `gen-…` when the display name normalizes to nothing.
 * Pure and stable across renders (unlike genRandomChars in translateDisplayNameForK8sAndReport),
 * so callers can safely derive auto K8s names in useMemo and compare to form state for touched.
 */
const stableGenSuffixFromSeed = (seed: string): string => {
  let h = 2166136261;
  for (let i = 0; i < seed.length; i++) {
    h ^= seed.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  let x = h >>> 0;
  let out = '';
  for (let i = 0; i < 4; i++) {
    x = Math.imul(x, 1664525) + 1013904223;
    out += (x >>> 20).toString(36);
  }
  const alphanumeric = out.replace(/[^a-z0-9]/g, '');
  return (alphanumeric.length >= 4 ? alphanumeric : `${alphanumeric}x0`).slice(0, 16);
};

/**
 * Translates a name to a k8s-safe value.
 * When the display name is non-empty but normalizes to nothing (e.g. dashes or punctuation only),
 * returns `gen-` + a stable suffix derived from the trimmed input (same intent as
 * translateDisplayNameForK8sAndReport in app/shared/components/utils.ts, without per-call randomness).
 * @see https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
 */
export const translateDisplayNameForK8s = (name = '', safePrefix = ''): string => {
  const trimmedInput = name.trim();
  const translatedName = trimmedInput
    .toLowerCase()
    .replace(/\s/g, '-')
    .replace(/[^a-z0-9-]/g, '')
    .replace(/^-*/, '')
    .replace(/-*$/, '')
    .replace(/[-]+/g, '-');

  if (safePrefix) {
    return `${safePrefix}${translatedName}`;
  }

  if (trimmedInput.length > 0 && translatedName.trim().length === 0) {
    return `gen-${stableGenSuffixFromSeed(trimmedInput)}`;
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
  const { initialData, editableK8sName, maxK8sNameLength } = configuration;
  const maxLength = maxK8sNameLength ?? MAX_K8S_NAME_LENGTH;

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
        invalidLength: k8sName.length > maxLength,
        maxLength,
        touched: !!k8sName,
      },
    },
  };
};

export const handleUpdateLogic =
  (currentData: K8sNameDescriptionFieldData): K8sNameDescriptionFieldUpdateFunctionInternal =>
  (key, value) => {
    const { maxLength } = currentData.k8sName.state;

    if (key === 'k8sName') {
      const validation = checkValidK8sName(value);
      return {
        ...currentData,
        k8sName: {
          value,
          state: {
            ...currentData.k8sName.state,
            invalidCharacters: !validation.valid && validation.invalidCharacters,
            invalidLength: value.length > maxLength,
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
            invalidLength: k8sValue.length > maxLength,
          },
        },
      };
    }

    return {
      ...currentData,
      [key]: value,
    };
  };
