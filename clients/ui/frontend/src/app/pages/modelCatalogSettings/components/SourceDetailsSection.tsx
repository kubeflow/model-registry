import * as React from 'react';
import {
  FormGroup,
  TextInput,
  Radio,
  FormHelperText,
  HelperText,
  HelperTextItem,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
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
  setData: UpdateObjectAtPropAndValue<ManageSourceFormData>;
};

const SourceDetailsSection: React.FC<SourceDetailsSectionProps> = ({ formData, setData }) => {
  const [isNameTouched, setIsNameTouched] = React.useState(false);
  const isNameValid = validateSourceName(formData.name);

  const nameInput = (
    <TextInput
      isRequired
      type="text"
      id="source-name"
      name="source-name"
      data-testid="source-name-input"
      value={formData.name}
      onChange={(_event, value) => setData('name', value)}
      onBlur={() => setIsNameTouched(true)}
      validated={isNameTouched && !isNameValid ? 'error' : 'default'}
    />
  );

  return (
    <FormSection>
      <FormGroup label={FORM_LABELS.NAME} isRequired fieldId="source-name">
        <FormFieldset component={nameInput} field="Name" />
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
              onChange={() => setData('sourceType', SourceType.HuggingFace)}
              label={SOURCE_TYPE_LABELS.HUGGING_FACE}
              id="source-type-huggingface"
              data-testid="source-type-huggingface"
            />
          </FlexItem>
          <FlexItem>
            <Radio
              isChecked={formData.sourceType === SourceType.YAML}
              name="source-type"
              onChange={() => setData('sourceType', SourceType.YAML)}
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
