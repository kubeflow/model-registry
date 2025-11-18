import * as React from 'react';
import {
  FormGroup,
  TextArea,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { validateYamlContent } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  VALIDATION_MESSAGES,
  HELP_TEXT,
} from '~/app/pages/modelCatalogSettings/constants';

type YamlSectionProps = {
  formData: ManageSourceFormData;
  setData: UpdateObjectAtPropAndValue<ManageSourceFormData>;
};

const YamlSection: React.FC<YamlSectionProps> = ({ formData, setData }) => {
  const [isYamlTouched, setIsYamlTouched] = React.useState(false);
  const isYamlContentValid = validateYamlContent(formData.yamlContent);

  const yamlInput = (
    <TextArea
      isRequired
      id="yaml-content"
      name="yaml-content"
      data-testid="yaml-content-input"
      value={formData.yamlContent}
      onChange={(_event, value) => setData('yamlContent', value)}
      onBlur={() => setIsYamlTouched(true)}
      validated={isYamlTouched && !isYamlContentValid ? 'error' : 'default'}
      rows={10}
      resizeOrientation="vertical"
    />
  );

  return (
    <FormSection data-testid="yaml-section">
      <FormGroup label={FORM_LABELS.YAML_CONTENT} isRequired fieldId="yaml-content">
        <FormFieldset component={yamlInput} field="YAML" />
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{HELP_TEXT.YAML}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        {isYamlTouched && !isYamlContentValid && (
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
