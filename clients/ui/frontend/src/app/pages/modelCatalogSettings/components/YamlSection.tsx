import * as React from 'react';
import {
  FormSection,
  FormGroup,
  TextArea,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { validateYamlContent } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  VALIDATION_MESSAGES,
  HELP_TEXT,
} from '~/app/pages/modelCatalogSettings/constants';

type YamlSectionProps = {
  formData: ManageSourceFormData;
  touched: Record<string, boolean>;
  onDataChange: (key: keyof ManageSourceFormData, value: string) => void;
  onFieldBlur: (field: string) => void;
};

const YamlSection: React.FC<YamlSectionProps> = ({
  formData,
  touched,
  onDataChange,
  onFieldBlur,
}) => {
  const isYamlContentValid = validateYamlContent(formData.yamlContent);

  return (
    <FormSection data-testid="yaml-section">
      <FormGroup label={FORM_LABELS.YAML_CONTENT} isRequired fieldId="yaml-content">
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{HELP_TEXT.YAML}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        <TextArea
          isRequired
          id="yaml-content"
          name="yaml-content"
          data-testid="yaml-content-input"
          value={formData.yamlContent}
          onChange={(_event, value) => onDataChange('yamlContent', value)}
          onBlur={() => onFieldBlur('yamlContent')}
          validated={touched.yamlContent && !isYamlContentValid ? 'error' : 'default'}
          rows={10}
          resizeOrientation="vertical"
        />
        {touched.yamlContent && !isYamlContentValid && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem variant="error" data-testid="yaml-content-error">
                {VALIDATION_MESSAGES.YAML_CONTENT_REQUIRED}
              </HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
      </FormGroup>
    </FormSection>
  );
};

export default YamlSection;
