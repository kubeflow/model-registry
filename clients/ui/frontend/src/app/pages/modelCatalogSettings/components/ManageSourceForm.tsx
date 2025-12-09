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
  useManageSourceData,
} from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { FORM_LABELS, DESCRIPTIONS } from '~/app/pages/modelCatalogSettings/constants';
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import { transformFormDataToPayload } from '~/app/pages/modelCatalogSettings/utils/modelCatalogSettingsUtils';
import {
  CatalogSourceType,
  CatalogSourcePreviewResult,
  CatalogSourcePreviewRequest,
} from '~/app/modelCatalogTypes';
import SourceDetailsSection from './SourceDetailsSection';
import CredentialsSection from './CredentialsSection';
import YamlSection from './YamlSection';
import ModelVisibilitySection from './ModelVisibilitySection';
import PreviewPanel from './PreviewPanel';
import ManageSourceFormFooter from './ManageSourceFormFooter';

type PreviewState = {
  mode?: 'preview' | 'validate';
  isLoading: boolean;
  result?: CatalogSourcePreviewResult;
  error?: Error;
  resultDismissed: boolean;
  lastPreviewedData?: string;
};

type ManageSourceFormProps = {
  existingData?: Partial<ManageSourceFormData>;
  isEditMode: boolean;
};

const ManageSourceForm: React.FC<ManageSourceFormProps> = ({ existingData, isEditMode }) => {
  const navigate = useNavigate();
  const [formData, setData] = useManageSourceData(existingData);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);
  const { apiState, refreshCatalogSourceConfigs } = React.useContext(ModelCatalogSettingsContext);

  // Preview state
  const [previewState, setPreviewState] = React.useState<PreviewState>({
    isLoading: false,
    resultDismissed: false,
  });

  const isHuggingFaceMode = formData.sourceType === CatalogSourceType.HUGGING_FACE;
  const isFormComplete = isFormValid(formData);
  const canPreview = isPreviewReady(formData);

  // Derive whether form has changed since last preview
  const hasFormChanged = React.useMemo(() => {
    const currentData = JSON.stringify(formData);
    return (
      previewState.lastPreviewedData !== undefined && currentData !== previewState.lastPreviewedData
    );
  }, [formData, previewState.lastPreviewedData]);

  // Derive validation success state
  const isValidationSuccess =
    previewState.mode === 'validate' &&
    !previewState.isLoading &&
    !previewState.error &&
    !previewState.resultDismissed;

  // Auto-trigger preview on mount in edit mode
  React.useEffect(() => {
    if (isEditMode && existingData && canPreview && !previewState.result) {
      handlePreview();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const buildPreviewRequest = (): CatalogSourcePreviewRequest => {
    const payload = transformFormDataToPayload(formData);

    const request: CatalogSourcePreviewRequest = {
      type: payload.type,
      includedModels: payload.includedModels,
      excludedModels: payload.excludedModels,
    };

    if (payload.type === CatalogSourceType.HUGGING_FACE) {
      request.properties = {
        allowedOrganization: payload.allowedOrganization,
        apiKey: payload.apiKey,
      };
    } else {
      request.properties = {
        yaml: payload.yaml,
      };
    }

    return request;
  };

  const handlePreview = async (mode: 'preview' | 'validate' = 'preview') => {
    if (!apiState.apiAvailable) {
      setPreviewState({
        mode,
        isLoading: false,
        error: new Error('API is not available'),
        resultDismissed: false,
      });
      return;
    }

    // Start loading, clear previous state
    setPreviewState({
      mode,
      isLoading: true,
      resultDismissed: false,
    });

    try {
      const previewRequest = buildPreviewRequest();
      const result = await apiState.api.previewCatalogSource({}, previewRequest);

      setPreviewState({
        mode,
        isLoading: false,
        result,
        lastPreviewedData: JSON.stringify(formData),
        resultDismissed: false,
      });
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Failed to preview source');

      setPreviewState({
        mode,
        isLoading: false,
        error: err,
        resultDismissed: false,
      });
    }
  };

  const handleValidate = async () => {
    await handlePreview('validate');
  };

  const handleSubmit = async () => {
    if (!apiState.apiAvailable) {
      setSubmitError(new Error('API is not available'));
      return;
    }
    setIsSubmitting(true);
    setSubmitError(undefined);

    try {
      const payload = transformFormDataToPayload(formData);

      if (isEditMode) {
        await apiState.api.updateCatalogSourceConfig({}, payload.id, payload);
      } else {
        await apiState.api.createCatalogSourceConfig({}, payload);
      }

      refreshCatalogSourceConfigs();
      navigate(catalogSettingsUrl());
    } catch (error) {
      setSubmitError(error instanceof Error ? error : new Error(`Failed to save source`));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    navigate(catalogSettingsUrl());
  };

  return (
    <>
      <Sidebar hasBorder isPanelRight hasGutter>
        <SidebarContent>
          <Form isWidthLimited>
            <Stack hasGutter>
              <StackItem>
                <SourceDetailsSection
                  formData={formData}
                  setData={setData}
                  isEditMode={isEditMode}
                />
              </StackItem>

              {isHuggingFaceMode && (
                <StackItem>
                  <CredentialsSection
                    formData={formData}
                    setData={setData}
                    onValidate={handleValidate}
                    isValidating={previewState.mode === 'validate' && previewState.isLoading}
                    validationError={
                      previewState.mode === 'validate' ? previewState.error : undefined
                    }
                    isValidationSuccess={isValidationSuccess}
                    onClearValidationSuccess={() =>
                      setPreviewState({ ...previewState, resultDismissed: true })
                    }
                  />
                </StackItem>
              )}

              {!formData.isDefault && !isHuggingFaceMode && (
                <StackItem>
                  <YamlSection formData={formData} setData={setData} />
                </StackItem>
              )}

              <StackItem>
                <ModelVisibilitySection
                  formData={formData}
                  setData={setData}
                  isDefaultExpanded={
                    existingData?.isDefault ||
                    !!existingData?.allowedModels ||
                    !!existingData?.excludedModels
                  }
                />
              </StackItem>

              <StackItem>
                <FormSection>
                  <FormGroup fieldId="enable-source">
                    <Checkbox
                      label={
                        <span className="pf-v6-c-form__label-text">
                          {FORM_LABELS.ENABLE_SOURCE}
                        </span>
                      }
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
        <SidebarPanel width={{ default: 'width_50' }}>
          <PreviewPanel
            isPreviewEnabled={canPreview}
            isLoading={previewState.isLoading}
            onPreview={() => handlePreview('preview')}
            previewResult={previewState.result}
            previewError={previewState.mode === 'preview' ? previewState.error : undefined}
            hasFormChanged={hasFormChanged}
          />
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
        isPreviewLoading={previewState.isLoading}
        onPreview={() => handlePreview('preview')}
      />
    </>
  );
};

export default ManageSourceForm;
