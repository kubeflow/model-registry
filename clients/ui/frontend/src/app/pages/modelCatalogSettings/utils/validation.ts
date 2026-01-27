import { CatalogSourceType } from '~/app/modelCatalogTypes';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { SOURCE_NAME_CHARACTER_LIMIT } from '~/app/pages/modelCatalogSettings/constants';

const isNonEmptyString = (value: string): boolean => value.trim().length > 0;

export const validateSourceName = (name: string): boolean =>
  isNonEmptyString(name) && name.length <= SOURCE_NAME_CHARACTER_LIMIT;

export const isSourceNameEmpty = (name: string): boolean => !isNonEmptyString(name);

export const validateOrganization = (organization: string): boolean =>
  isNonEmptyString(organization);

export const validateYamlContent = (yamlContent: string): boolean => isNonEmptyString(yamlContent);

export const validateHuggingFaceCredentials = (data: ManageSourceFormData): boolean => {
  if (data.sourceType !== CatalogSourceType.HUGGING_FACE) {
    return true;
  }
  return validateOrganization(data.organization);
};

export const validateYamlMode = (data: ManageSourceFormData): boolean => {
  if (data.sourceType !== CatalogSourceType.YAML || data.isDefault) {
    return true;
  }
  return validateYamlContent(data.yamlContent);
};

export const isFormValid = (data: ManageSourceFormData): boolean =>
  validateSourceName(data.name) && validateHuggingFaceCredentials(data) && validateYamlMode(data);

export const isPreviewReady = (data: ManageSourceFormData): boolean => {
  if (data.sourceType === CatalogSourceType.HUGGING_FACE) {
    return validateHuggingFaceCredentials(data);
  }
  return validateYamlMode(data);
};
