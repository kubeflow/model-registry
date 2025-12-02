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

export const transformFormDataToPayload = (
  formData: ManageSourceFormData,
): CatalogSourceConfigPayload => ({
  id: formData.id || generateSourceIdFromName(formData.name),
  name: formData.name,
  type: formData.sourceType,
  enabled: formData.enabled,
  isDefault: formData.isDefault,
  includedModels: formData.allowedModels
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0),
  excludedModels: formData.excludedModels
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0),
  ...(formData.sourceType === CatalogSourceType.YAML && { yaml: formData.yamlContent }),
  ...(formData.sourceType === CatalogSourceType.HUGGING_FACE && {
    apiKey: formData.accessToken,
    allowedOrganization: formData.organization,
  }),
});
