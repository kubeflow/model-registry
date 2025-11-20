import { GenericObjectState } from 'mod-arch-core';
import useGenericObjectState from 'mod-arch-core/dist/utilities/useGenericObjectState';

export enum SourceType {
  HuggingFace = 'huggingface',
  YAML = 'yaml',
}

export type ManageSourceFormData = {
  name: string;
  sourceType: SourceType;
  // Hugging Face fields
  accessToken: string;
  organization: string;
  // YAML field
  yamlContent: string;
  // Filter fields
  allowedModels: string;
  excludedModels: string;
  // Enable source
  enabled: boolean;
};

const manageSourceFormDataDefaults: ManageSourceFormData = {
  name: '',
  sourceType: SourceType.HuggingFace,
  accessToken: '',
  organization: '',
  yamlContent: '',
  allowedModels: '',
  excludedModels: '',
  enabled: false,
};

/**
 * Custom hook to manage form state for adding/editing a catalog source
 * Uses the standard useGenericObjectState pattern from mod-arch-core
 * @param existingData - Optional existing data to pre-populate the form (for edit mode)
 * @returns Generic object state with [formData, setData] pattern
 */
export const useManageSourceData = (
  existingData?: Partial<ManageSourceFormData>,
): GenericObjectState<ManageSourceFormData> =>
  useGenericObjectState<ManageSourceFormData>({
    ...manageSourceFormDataDefaults,
    ...existingData,
  });
