import * as React from 'react';
import {
  ExpandableSection,
  FormSection,
  FormGroup,
  TextArea,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import {
  ManageSourceFormData,
  SourceType,
} from '~/app/pages/modelCatalogSettings/useManageSourceData';
import {
  FORM_LABELS,
  PLACEHOLDERS,
  DESCRIPTIONS,
  FIELD_HELPER_TEXT,
  getFilterInfoWithOrg,
  getAllowedModelsHelp,
  getExcludedModelsHelp,
} from '~/app/pages/modelCatalogSettings/constants';

type ModelVisibilitySectionProps = {
  formData: ManageSourceFormData;
  isExpanded: boolean;
  onToggle: () => void;
  onDataChange: (key: keyof ManageSourceFormData, value: string) => void;
};

const ModelVisibilitySection: React.FC<ModelVisibilitySectionProps> = ({
  formData,
  isExpanded,
  onToggle,
  onDataChange,
}) => {
  const isHuggingFaceMode = formData.sourceType === SourceType.HuggingFace;
  const organization = isHuggingFaceMode ? formData.organization : undefined;

  const sectionDescription =
    isHuggingFaceMode && organization
      ? getFilterInfoWithOrg(organization)
      : DESCRIPTIONS.FILTER_INFO_GENERIC;

  const allowedModelsHelp = getAllowedModelsHelp(organization);
  const excludedModelsHelp = getExcludedModelsHelp(organization);

  const allowedModelsPlaceholder = isHuggingFaceMode
    ? PLACEHOLDERS.ALLOWED_MODELS_HF
    : PLACEHOLDERS.ALLOWED_MODELS_GENERIC;

  const excludedModelsPlaceholder = isHuggingFaceMode
    ? PLACEHOLDERS.EXCLUDED_MODELS_HF
    : PLACEHOLDERS.EXCLUDED_MODELS_GENERIC;

  return (
    <ExpandableSection
      toggleText={FORM_LABELS.MODEL_VISIBILITY}
      onToggle={onToggle}
      isExpanded={isExpanded}
      data-testid="model-visibility-section"
    >
      <FormSection>
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{sectionDescription}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        <FormGroup label={FORM_LABELS.ALLOWED_MODELS} fieldId="allowed-models">
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{allowedModelsHelp}</HelperTextItem>
            </HelperText>
          </FormHelperText>
          <TextArea
            id="allowed-models"
            name="allowed-models"
            data-testid="allowed-models-input"
            value={formData.allowedModels}
            onChange={(_event, value) => onDataChange('allowedModels', value)}
            rows={3}
            resizeOrientation="vertical"
            placeholder={allowedModelsPlaceholder}
          />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{FIELD_HELPER_TEXT.INCLUDED_MODELS}</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>

        <FormGroup label={FORM_LABELS.EXCLUDED_MODELS} fieldId="excluded-models">
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{excludedModelsHelp}</HelperTextItem>
            </HelperText>
          </FormHelperText>
          <TextArea
            id="excluded-models"
            name="excluded-models"
            data-testid="excluded-models-input"
            value={formData.excludedModels}
            onChange={(_event, value) => onDataChange('excludedModels', value)}
            rows={3}
            resizeOrientation="vertical"
            placeholder={excludedModelsPlaceholder}
          />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{FIELD_HELPER_TEXT.EXCLUDED_MODELS}</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>
      </FormSection>
    </ExpandableSection>
  );
};

export default ModelVisibilitySection;
