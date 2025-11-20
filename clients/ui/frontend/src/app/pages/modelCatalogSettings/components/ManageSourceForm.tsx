import * as React from 'react';
import {
  Form,
  FormGroup,
  Checkbox,
  Stack,
  StackItem,
  Sidebar,
  SidebarPanel,
  SidebarContent,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { catalogSettingsUrl } from '~/app/routes/modelCatalogSettings/modelCatalogSettings';
import { isFormValid, isPreviewReady } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  ManageSourceFormData,
  SourceType,
  useManageSourceData,
} from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { FORM_LABELS, DESCRIPTIONS } from '~/app/pages/modelCatalogSettings/constants';
import SourceDetailsSection from './SourceDetailsSection';
import CredentialsSection from './CredentialsSection';
import YamlSection from './YamlSection';
import ModelVisibilitySection from './ModelVisibilitySection';
import PreviewPanel from './PreviewPanel';
import ManageSourceFormFooter from './ManageSourceFormFooter';

type ManageSourceFormProps = {
  existingData?: Partial<ManageSourceFormData>;
  isEditMode: boolean;
};

const ManageSourceForm: React.FC<ManageSourceFormProps> = ({ existingData, isEditMode }) => {
  const navigate = useNavigate();
  const [formData, setData] = useManageSourceData(existingData);
  const [isFiltersExpanded, setIsFiltersExpanded] = React.useState(false);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);

  const isHuggingFaceMode = formData.sourceType === SourceType.HuggingFace;
  const isFormComplete = isFormValid(formData);
  const canPreview = isPreviewReady(formData);

  const handleSubmit = async () => {
    setIsSubmitting(true);
    setSubmitError(undefined);

    try {
      // TODO: Implement submit logic (will be part of API integration)
      // navigate(catalogSettingsUrl());
    } catch (error) {
      setSubmitError(error instanceof Error ? error : new Error('Failed to save source'));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePreview = () => {
    // TODO: Implement preview logic (will be part of API integration)
  };

  const handleCancel = () => {
    navigate(catalogSettingsUrl());
  };

  return (
    <>
      <Sidebar hasGutter hasBorder isPanelRight>
        <SidebarContent>
          <Form isWidthLimited>
            <Stack hasGutter>
              <StackItem>
                <SourceDetailsSection formData={formData} setData={setData} />
              </StackItem>

              {isHuggingFaceMode && (
                <StackItem>
                  <CredentialsSection formData={formData} setData={setData} />
                </StackItem>
              )}

              {!isHuggingFaceMode && (
                <StackItem>
                  <YamlSection formData={formData} setData={setData} />
                </StackItem>
              )}

              <StackItem>
                <ModelVisibilitySection
                  formData={formData}
                  isExpanded={isFiltersExpanded}
                  onToggle={() => setIsFiltersExpanded((prev) => !prev)}
                  setData={setData}
                />
              </StackItem>

              <StackItem>
                <FormSection>
                  <FormGroup fieldId="enable-source">
                    <Checkbox
                      label={<strong>{FORM_LABELS.ENABLE_SOURCE}</strong>}
                      id="enable-source"
                      name="enable-source"
                      data-testid="enable-source-checkbox"
                      description={DESCRIPTIONS.ENABLE_SOURCE}
                      isChecked={formData.enabled}
                      onChange={(_event, checked) => setData('enabled', checked)}
                    />
                  </FormGroup>
                </FormSection>
              </StackItem>
            </Stack>
          </Form>
        </SidebarContent>
        <SidebarPanel hasPadding className="pf-v6-u-flex-basis-40">
          <PreviewPanel isPreviewEnabled={canPreview} onPreview={handlePreview} />
        </SidebarPanel>
      </Sidebar>
      <ManageSourceFormFooter
        submitLabel={isEditMode ? 'Save' : 'Add'}
        submitError={submitError}
        isSubmitDisabled={!isFormComplete || isSubmitting}
        isSubmitting={isSubmitting}
        onSubmit={handleSubmit}
        onCancel={handleCancel}
        isPreviewDisabled={!canPreview}
        onPreview={handlePreview}
      />
    </>
  );
};

export default ManageSourceForm;
