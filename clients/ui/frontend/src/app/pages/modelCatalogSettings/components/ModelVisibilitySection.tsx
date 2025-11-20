import * as React from 'react';
import {
  FormFieldGroupExpandable,
  FormFieldGroupHeader,
  FormGroup,
  TextArea,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
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
  setData: UpdateObjectAtPropAndValue<ManageSourceFormData>;
};

const ModelVisibilitySection: React.FC<ModelVisibilitySectionProps> = ({ formData, setData }) => {
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

  const allowedModelsInput = (
    <TextArea
      id="allowed-models"
      name="allowed-models"
      data-testid="allowed-models-input"
      value={formData.allowedModels}
      onChange={(_event, value) => setData('allowedModels', value)}
      rows={3}
      resizeOrientation="vertical"
      placeholder={allowedModelsPlaceholder}
    />
  );

  const excludedModelsInput = (
    <TextArea
      id="excluded-models"
      name="excluded-models"
      data-testid="excluded-models-input"
      value={formData.excludedModels}
      onChange={(_event, value) => setData('excludedModels', value)}
      rows={3}
      resizeOrientation="vertical"
      placeholder={excludedModelsPlaceholder}
    />
  );

  return (
    <FormSection>
      <FormFieldGroupExpandable
        toggleAriaLabel="Model visibility"
        header={
          <FormFieldGroupHeader
            titleText={{ text: FORM_LABELS.MODEL_VISIBILITY, id: 'model-visibility-title' }}
            titleDescription={sectionDescription}
          />
        }
        data-testid="model-visibility-section"
      >
        <FormGroup label={FORM_LABELS.ALLOWED_MODELS} fieldId="allowed-models">
          <FormFieldset component={allowedModelsInput} field="Allowed models" />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{allowedModelsHelp}</HelperTextItem>
            </HelperText>
          </FormHelperText>
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{FIELD_HELPER_TEXT.INCLUDED_MODELS}</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>

        <FormGroup label={FORM_LABELS.EXCLUDED_MODELS} fieldId="excluded-models">
          <FormFieldset component={excludedModelsInput} field="Excluded models" />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{excludedModelsHelp}</HelperTextItem>
            </HelperText>
          </FormHelperText>
          <FormHelperText>
            <HelperText>
              <HelperTextItem>{FIELD_HELPER_TEXT.EXCLUDED_MODELS}</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>
      </FormFieldGroupExpandable>
    </FormSection>
  );
};

export default ModelVisibilitySection;
