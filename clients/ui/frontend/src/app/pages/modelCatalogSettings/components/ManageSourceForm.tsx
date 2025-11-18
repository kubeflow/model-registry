import * as React from 'react';
import {
  Form,
  FormSection,
  FormGroup,
  Checkbox,
  Button,
  ActionGroup,
  Grid,
  GridItem,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { catalogSettingsUrl } from '~/app/routes/modelCatalogSettings/modelCatalogSettings';
import { isFormValid, isPreviewReady } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  ManageSourceFormData,
  SourceType,
  useManageSourceData,
} from '~/app/pages/modelCatalogSettings/useManageSourceData';
import {
  FORM_LABELS,
  BUTTON_LABELS,
  DESCRIPTIONS,
} from '~/app/pages/modelCatalogSettings/constants';
import SourceDetailsSection from './SourceDetailsSection';
import CredentialsSection from './CredentialsSection';
import YamlSection from './YamlSection';
import ModelVisibilitySection from './ModelVisibilitySection';
import PreviewButton from './PreviewButton';
import PreviewPanel from './PreviewPanel';

type ManageSourceFormProps = {
  existingData?: Partial<ManageSourceFormData>;
  isEditMode: boolean;
};

const ManageSourceForm: React.FC<ManageSourceFormProps> = ({ existingData, isEditMode }) => {
  const navigate = useNavigate();
  const { formData, touched, updateField, markFieldAsTouched } = useManageSourceData(existingData);
  const [isFiltersExpanded, setIsFiltersExpanded] = React.useState(false);

  const isHuggingFaceMode = formData.sourceType === SourceType.HuggingFace;
  const isFormComplete = isFormValid(formData);
  const canPreview = isPreviewReady(formData);

  const handleSubmit = () => {
    // TODO: Implement submit logic (will be part of API integration)
  };

  const handlePreview = () => {
    // TODO: Implement preview logic (will be part of API integration)
  };

  const handleCancel = () => {
    navigate(catalogSettingsUrl());
  };

  const handleToggleFilters = () => {
    setIsFiltersExpanded((prev) => !prev);
  };

  return (
    <Grid hasGutter span={12} style={{ height: '100%' }}>
      <GridItem span={7}>
        <Form>
          <SourceDetailsSection
            formData={formData}
            touched={touched}
            onDataChange={updateField}
            onFieldBlur={markFieldAsTouched}
          />

          {isHuggingFaceMode && (
            <CredentialsSection
              formData={formData}
              touched={touched}
              onDataChange={updateField}
              onFieldBlur={markFieldAsTouched}
            />
          )}

          {!isHuggingFaceMode && (
            <YamlSection
              formData={formData}
              touched={touched}
              onDataChange={updateField}
              onFieldBlur={markFieldAsTouched}
            />
          )}

          <ModelVisibilitySection
            formData={formData}
            isExpanded={isFiltersExpanded}
            onToggle={handleToggleFilters}
            onDataChange={updateField}
          />

          <FormSection>
            <FormGroup fieldId="enable-source">
              <Checkbox
                label={FORM_LABELS.ENABLE_SOURCE}
                id="enable-source"
                name="enable-source"
                data-testid="enable-source-checkbox"
                description={DESCRIPTIONS.ENABLE_SOURCE}
                isChecked={formData.enabled}
                onChange={(_event, checked) => updateField('enabled', checked)}
              />
            </FormGroup>
          </FormSection>

          <ActionGroup>
            <Button
              variant="primary"
              onClick={handleSubmit}
              isDisabled={!isFormComplete}
              data-testid="submit-button"
            >
              {isEditMode ? BUTTON_LABELS.SAVE : BUTTON_LABELS.ADD}
            </Button>
            <PreviewButton
              onClick={handlePreview}
              isDisabled={!canPreview}
              variant="secondary"
              testId="preview-button"
            />
            <Button variant="link" onClick={handleCancel} data-testid="cancel-button">
              {BUTTON_LABELS.CANCEL}
            </Button>
          </ActionGroup>
        </Form>
      </GridItem>

      <GridItem span={5}>
        <PreviewPanel isPreviewEnabled={canPreview} onPreview={handlePreview} />
      </GridItem>
    </Grid>
  );
};

export default ManageSourceForm;
