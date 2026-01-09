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
import { useManageSourceData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { FORM_LABELS, DESCRIPTIONS } from '~/app/pages/modelCatalogSettings/constants';
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import {
  catalogSourceConfigToFormData,
  getPayloadForConfig,
  transformFormDataToConfig,
} from '~/app/pages/modelCatalogSettings/utils/modelCatalogSettingsUtils';
import {
  CatalogSourceConfig,
  CatalogSourceType,
  CatalogSourcePreviewRequest,
  CatalogSourcePreviewModel,
  CatalogSourcePreviewSummary,
} from '~/app/modelCatalogTypes';
import SourceDetailsSection from './SourceDetailsSection';
import CredentialsSection from './CredentialsSection';
import YamlSection from './YamlSection';
import ModelVisibilitySection from './ModelVisibilitySection';
import PreviewPanel from './PreviewPanel';
import ManageSourceFormFooter from './ManageSourceFormFooter';

type FilterStatus = 'included' | 'excluded';

type PreviewTabState = {
  items: CatalogSourcePreviewModel[];
  nextPageToken?: string;
  hasMore: boolean;
};

const initialTabState: PreviewTabState = {
  items: [],
  nextPageToken: undefined,
  hasMore: false,
};

type PreviewState = {
  mode?: 'preview' | 'validate';
  isLoadingInitial: boolean;
  isLoadingMore: boolean;
  summary?: CatalogSourcePreviewSummary;
  tabStates: {
    included: PreviewTabState;
    excluded: PreviewTabState;
  };
  error?: Error;
  resultDismissed: boolean;
  lastPreviewedData?: CatalogSourcePreviewRequest;
  activeTab: FilterStatus;
};

type ManageSourceFormProps = {
  existingSourceConfig?: CatalogSourceConfig;
  isEditMode: boolean;
};

const ManageSourceForm: React.FC<ManageSourceFormProps> = ({
  existingSourceConfig,
  isEditMode,
}) => {
  const navigate = useNavigate();
  const existingData = existingSourceConfig
    ? catalogSourceConfigToFormData(existingSourceConfig)
    : undefined;
  const [formData, setData] = useManageSourceData(existingData);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);
  const { apiState, refreshCatalogSourceConfigs } = React.useContext(ModelCatalogSettingsContext);

  // Preview state
  const [previewState, setPreviewState] = React.useState<PreviewState>({
    isLoadingInitial: false,
    isLoadingMore: false,
    tabStates: {
      included: initialTabState,
      excluded: initialTabState,
    },
    resultDismissed: false,
    activeTab: 'included',
  });

  const isHuggingFaceMode = formData.sourceType === CatalogSourceType.HUGGING_FACE;
  const isFormComplete = isFormValid(formData);
  const canPreview = isPreviewReady(formData);

  const buildPreviewRequest = React.useCallback((): CatalogSourcePreviewRequest => {
    const payload = transformFormDataToConfig(formData, existingSourceConfig);

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
        yamlCatalogPath: payload.yamlCatalogPath,
      };
    }

    return request;
  }, [formData, existingSourceConfig]);

  // Derive whether form has changed since last preview
  const hasFormChanged = React.useMemo(() => {
    if (!previewState.lastPreviewedData) {
      return false;
    }
    const currentRequest = buildPreviewRequest();
    return JSON.stringify(currentRequest) !== JSON.stringify(previewState.lastPreviewedData);
  }, [buildPreviewRequest, previewState.lastPreviewedData]);

  // Derive validation success state
  const isValidationSuccess =
    previewState.mode === 'validate' &&
    !previewState.isLoadingInitial &&
    !previewState.error &&
    !previewState.resultDismissed;

  // Auto-trigger preview on mount in edit mode
  React.useEffect(() => {
    const hasNoResults = previewState.tabStates.included.items.length === 0;
    if (isEditMode && existingData && canPreview && hasNoResults) {
      handlePreview();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handlePreview = async (
    mode: 'preview' | 'validate' = 'preview',
    options?: {
      loadMore?: boolean;
      switchToTab?: FilterStatus;
    },
  ) => {
    const { loadMore = false, switchToTab } = options ?? {};
    const isFreshPreview = !loadMore && !switchToTab;
    // For fresh preview, always start with 'included' tab
    const targetTab = isFreshPreview ? 'included' : (switchToTab ?? previewState.activeTab);

    if (!apiState.apiAvailable) {
      setPreviewState((prev) => ({
        ...prev,
        mode,
        isLoadingInitial: false,
        error: new Error('API is not available'),
        resultDismissed: false,
      }));
      return;
    }

    // For fresh preview, reset everything to clean state
    if (isFreshPreview) {
      setPreviewState({
        mode,
        isLoadingInitial: true,
        isLoadingMore: false,
        tabStates: { included: initialTabState, excluded: initialTabState },
        activeTab: 'included',
        error: undefined,
        resultDismissed: false,
        summary: undefined,
        lastPreviewedData: undefined,
      });
    } else if (loadMore) {
      setPreviewState((prev) => ({ ...prev, isLoadingMore: true }));
    } else if (switchToTab) {
      setPreviewState((prev) => ({ ...prev, activeTab: switchToTab, isLoadingInitial: true }));
    }

    // Use lastPreviewedData for load more / tab switch, current formData for fresh
    let requestData: CatalogSourcePreviewRequest;
    if (isFreshPreview) {
      requestData = buildPreviewRequest();
    } else if (previewState.lastPreviewedData) {
      requestData = previewState.lastPreviewedData;
    } else {
      // For non-fresh requests, lastPreviewedData must exist (guard against edge case)
      // eslint-disable-next-line no-console
      console.warn(
        'Attempted load more / tab switch without lastPreviewedData, triggering fresh preview',
      );
      return handlePreview(mode);
    }

    // Get token for load more
    const nextPageToken = loadMore ? previewState.tabStates[targetTab].nextPageToken : undefined;

    try {
      const result = await apiState.api.previewCatalogSource({}, requestData, {
        filterStatus: targetTab,
        pageSize: 20,
        nextPageToken,
      });

      // Update state based on operation type
      setPreviewState((prev) => {
        const currentTabState = prev.tabStates[targetTab];
        const newItems = loadMore ? [...currentTabState.items, ...result.items] : result.items;

        return {
          ...prev,
          mode,
          isLoadingInitial: false,
          isLoadingMore: false,
          summary: result.summary,
          lastPreviewedData: isFreshPreview ? requestData : prev.lastPreviewedData,
          tabStates: {
            ...prev.tabStates,
            [targetTab]: {
              items: newItems,
              nextPageToken: result.nextPageToken,
              hasMore: !!result.nextPageToken && result.items.length > 0,
            },
          },
          error: undefined,
          resultDismissed: false,
        };
      });
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Failed to preview source');

      setPreviewState((prev) => ({
        ...prev,
        mode,
        isLoadingInitial: false,
        isLoadingMore: false,
        error: err,
        resultDismissed: false,
      }));
    }
  };

  const handleTabChange = (newTab: FilterStatus) => {
    if (newTab === previewState.activeTab) {
      return;
    }

    const tabState = previewState.tabStates[newTab];
    if (tabState.items.length === 0) {
      // Tab not yet loaded, fetch it
      handlePreview('preview', { switchToTab: newTab });
    } else {
      // Tab already has data, just switch
      setPreviewState((prev) => ({ ...prev, activeTab: newTab }));
    }
  };

  const handleLoadMore = () => {
    handlePreview('preview', { loadMore: true });
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
      const sourceConfig = transformFormDataToConfig(formData, existingSourceConfig);
      const payload = getPayloadForConfig(sourceConfig, isEditMode);

      if (isEditMode) {
        await apiState.api.updateCatalogSourceConfig({}, formData.id, payload);
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
                    isValidating={previewState.mode === 'validate' && previewState.isLoadingInitial}
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
            isLoadingInitial={previewState.isLoadingInitial}
            onPreview={() => handlePreview('preview')}
            previewError={previewState.mode === 'preview' ? previewState.error : undefined}
            hasFormChanged={hasFormChanged}
            activeTab={previewState.activeTab}
            items={previewState.tabStates[previewState.activeTab].items}
            summary={previewState.summary}
            hasMore={previewState.tabStates[previewState.activeTab].hasMore}
            isLoadingMore={previewState.isLoadingMore}
            onTabChange={handleTabChange}
            onLoadMore={handleLoadMore}
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
        isPreviewLoading={previewState.isLoadingInitial}
        onPreview={() => handlePreview('preview')}
      />
    </>
  );
};

export default ManageSourceForm;
