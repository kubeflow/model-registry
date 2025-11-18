import * as React from 'react';
import {
  FormSection,
  FormGroup,
  TextInput,
  Radio,
  FormHelperText,
  HelperText,
  HelperTextItem,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import {
  ManageSourceFormData,
  SourceType,
} from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { validateSourceName } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  SOURCE_TYPE_LABELS,
  VALIDATION_MESSAGES,
} from '~/app/pages/modelCatalogSettings/constants';

type SourceDetailsSectionProps = {
  formData: ManageSourceFormData;
  touched: Record<string, boolean>;
  onDataChange: (key: keyof ManageSourceFormData, value: string | boolean) => void;
  onFieldBlur: (field: string) => void;
};

const SourceDetailsSection: React.FC<SourceDetailsSectionProps> = ({
  formData,
  touched,
  onDataChange,
  onFieldBlur,
}) => {
  const isNameValid = validateSourceName(formData.name);
  const isNameTouched = touched.name;

  return (
    <FormSection>
      <FormGroup label={FORM_LABELS.NAME} isRequired fieldId="source-name">
        <TextInput
          isRequired
          type="text"
          id="source-name"
          name="source-name"
          data-testid="source-name-input"
          value={formData.name}
          onChange={(_event, value) => onDataChange('name', value)}
          onBlur={() => onFieldBlur('name')}
          validated={isNameTouched && !isNameValid ? 'error' : 'default'}
        />
        {isNameTouched && !isNameValid && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem variant="error" data-testid="source-name-error">
                {VALIDATION_MESSAGES.NAME_REQUIRED}
              </HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
      </FormGroup>

      <FormGroup
        label={FORM_LABELS.SOURCE_TYPE}
        isRequired
        fieldId="source-type"
        role="radiogroup"
        aria-labelledby="source-type-label"
      >
        <Flex spaceItems={{ default: 'spaceItemsMd' }}>
          <FlexItem>
            <Radio
              isChecked={formData.sourceType === SourceType.HuggingFace}
              name="source-type"
              onChange={() => onDataChange('sourceType', SourceType.HuggingFace)}
              label={SOURCE_TYPE_LABELS.HUGGING_FACE}
              id="source-type-huggingface"
              data-testid="source-type-huggingface"
            />
          </FlexItem>
          <FlexItem>
            <Radio
              isChecked={formData.sourceType === SourceType.YAML}
              name="source-type"
              onChange={() => onDataChange('sourceType', SourceType.YAML)}
              label={SOURCE_TYPE_LABELS.YAML}
              id="source-type-yaml"
              data-testid="source-type-yaml"
            />
          </FlexItem>
        </Flex>
      </FormGroup>
    </FormSection>
  );
};

export default SourceDetailsSection;
