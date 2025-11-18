import {
  ManageSourceFormData,
  SourceType,
} from '~/app/pages/modelCatalogSettings/useManageSourceData';

const isNonEmptyString = (value: string): boolean => value.trim().length > 0;

export const validateSourceName = (name: string): boolean => isNonEmptyString(name);

export const validateOrganization = (organization: string): boolean =>
  isNonEmptyString(organization);

export const validateAccessToken = (accessToken: string): boolean => isNonEmptyString(accessToken);

export const validateYamlContent = (yamlContent: string): boolean => isNonEmptyString(yamlContent);

export const validateHuggingFaceCredentials = (data: ManageSourceFormData): boolean => {
  if (data.sourceType !== SourceType.HuggingFace) {
    return true;
  }
  return validateOrganization(data.organization) && validateAccessToken(data.accessToken);
};

export const validateYamlMode = (data: ManageSourceFormData): boolean => {
  if (data.sourceType !== SourceType.YAML) {
    return true;
  }
  return validateYamlContent(data.yamlContent);
};

export const isFormValid = (data: ManageSourceFormData): boolean =>
  validateSourceName(data.name) && validateHuggingFaceCredentials(data) && validateYamlMode(data);

export const isPreviewReady = (data: ManageSourceFormData): boolean => {
  if (data.sourceType === SourceType.HuggingFace) {
    return validateHuggingFaceCredentials(data);
  }
  return validateYamlMode(data);
};
