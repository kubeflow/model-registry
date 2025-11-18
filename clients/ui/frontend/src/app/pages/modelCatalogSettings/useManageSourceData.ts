import * as React from 'react';

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

type UseManageSourceDataReturn = {
  formData: ManageSourceFormData;
  touched: Record<string, boolean>;
  updateField: (key: keyof ManageSourceFormData, value: string | boolean) => void;
  markFieldAsTouched: (field: string) => void;
};

/**
 * Custom hook to manage form state for adding/editing a catalog source
 * @param existingData - Optional existing data to pre-populate the form (for edit mode)
 * @returns Form data, touched fields tracker, and handlers for updates
 */
export const useManageSourceData = (
  existingData?: Partial<ManageSourceFormData>,
): UseManageSourceDataReturn => {
  const [formData, setFormData] = React.useState<ManageSourceFormData>({
    name: existingData?.name ?? '',
    sourceType: existingData?.sourceType ?? SourceType.HuggingFace,
    accessToken: existingData?.accessToken ?? '',
    organization: existingData?.organization ?? '',
    yamlContent: existingData?.yamlContent ?? '',
    allowedModels: existingData?.allowedModels ?? '',
    excludedModels: existingData?.excludedModels ?? '',
    enabled: existingData?.enabled ?? false,
  });

  const [touched, setTouched] = React.useState<Record<string, boolean>>({});

  const updateField = React.useCallback(
    (key: keyof ManageSourceFormData, value: string | boolean) => {
      setFormData((prevData) => ({
        ...prevData,
        [key]: value,
      }));
    },
    [],
  );

  const markFieldAsTouched = React.useCallback((field: string) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
  }, []);

  return {
    formData,
    touched,
    updateField,
    markFieldAsTouched,
  };
};
