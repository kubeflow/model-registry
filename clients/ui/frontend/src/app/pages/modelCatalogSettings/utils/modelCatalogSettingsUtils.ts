import {
  CatalogSourceConfig,
  CatalogSourceConfigPayload,
  CatalogSourceType,
} from '~/app/modelCatalogTypes';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';

export const catalogSourceConfigToFormData = (
  sourceConfig: CatalogSourceConfig,
): Partial<ManageSourceFormData> => {
  const common: Partial<ManageSourceFormData> = {
    name: sourceConfig.name,
    sourceType: sourceConfig.type,
    enabled: sourceConfig.enabled ?? true,
    allowedModels: (sourceConfig.includedModels || []).join(', '),
    excludedModels: (sourceConfig.excludedModels || []).join(', '),
    isDefault: sourceConfig.isDefault,
    id: sourceConfig.id,
  };

  if (sourceConfig.type === CatalogSourceType.YAML) {
    return {
      ...common,
      yamlContent: sourceConfig.yaml ?? '',
      accessToken: '',
      organization: '',
    };
  }

  return {
    ...common,
    accessToken: sourceConfig.apiKey ?? '',
    organization: sourceConfig.allowedOrganization ?? '',
    yamlContent: '',
  };
};

export const generateSourceIdFromName = (name: string): string =>
  name
    .trim()
    .replace(/\s+/g, '_')
    .replace(/-/g, '_')
    .replace(/[^a-zA-Z0-9_]/g, '')
    .toLowerCase();

export const transformFormDataToConfig = (formData: ManageSourceFormData): CatalogSourceConfig => {
  const parseModels = (models: string): string[] =>
    models
      .split(',')
      .map((item) => item.trim())
      .filter((item) => item.length > 0);

  const commonFields = {
    id: formData.id || generateSourceIdFromName(formData.name),
    name: formData.name,
    enabled: formData.enabled,
    isDefault: formData.isDefault,
    includedModels: parseModels(formData.allowedModels),
    excludedModels: parseModels(formData.excludedModels),
  };

  if (formData.sourceType === CatalogSourceType.YAML) {
    return {
      ...commonFields,
      type: CatalogSourceType.YAML,
      yaml: formData.yamlContent,
    };
  }

  return {
    ...commonFields,
    type: CatalogSourceType.HUGGING_FACE,
    apiKey: formData.accessToken,
    allowedOrganization: formData.organization,
  };
};

export const getPayloadForConfig = (
  sourceConfig: CatalogSourceConfig,
  isEditMode = false,
): CatalogSourceConfigPayload => {
  if (sourceConfig.isDefault) {
    return {
      enabled: sourceConfig.enabled,
      includedModels: sourceConfig.includedModels,
      excludedModels: sourceConfig.excludedModels,
    };
  }

  if (isEditMode) {
    return {
      name: sourceConfig.name,
      type: sourceConfig.type,
      enabled: sourceConfig.enabled,
      isDefault: sourceConfig.isDefault,
      includedModels: sourceConfig.includedModels,
      excludedModels: sourceConfig.excludedModels,
      ...(sourceConfig.type === CatalogSourceType.YAML && { yaml: sourceConfig.yaml }),
      ...(sourceConfig.type === CatalogSourceType.HUGGING_FACE && {
        apiKey: sourceConfig.apiKey,
        allowedOrganization: sourceConfig.allowedOrganization,
      }),
    };
  }

  return sourceConfig;
};
