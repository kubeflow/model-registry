import {
  FormGroup,
  FormHelperText,
  HelperText,
  HelperTextItem,
  TextArea,
  TextInput,
} from '@patternfly/react-core';
import React from 'react';
import { UpdateObjectAtPropAndValue, FormSection } from 'mod-arch-shared';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { MR_CHARACTER_LIMIT } from './const';
import { RegisterModelFormData } from './useRegisterModelData';

type RegisterModelDetailsFormSectionProp<D extends RegisterModelFormData> = {
  formData: D;
  setData: UpdateObjectAtPropAndValue<D>;
  hasModelNameError: boolean;
  isModelNameDuplicate?: boolean;
};
const RegisterModelDetailsFormSection = <D extends RegisterModelFormData>({
  formData,
  setData,
  hasModelNameError,
  isModelNameDuplicate,
}: RegisterModelDetailsFormSectionProp<D>): React.ReactNode => {
  const modelNameInput = (
    <TextInput
      isRequired
      type="text"
      id="model-name"
      name="model-name"
      value={formData.modelName}
      onChange={(_e, value) => setData('modelName', value)}
      validated={hasModelNameError ? 'error' : 'default'}
    />
  );

  const modelDescriptionInput = (
    <TextArea
      type="text"
      id="model-description"
      name="model-description"
      value={formData.modelDescription}
      onChange={(_e, value) => setData('modelDescription', value)}
    />
  );

  return (
    <FormSection
      title="Model details"
      description="Provide model details that apply to every version of this model."
    >
      <FormGroup label="Model name" isRequired fieldId="model-name">
        <FormFieldset component={modelNameInput} />
        {hasModelNameError && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem variant="error" data-testid="model-name-error">
                {isModelNameDuplicate
                  ? 'Model name already exists'
                  : `Cannot exceed ${MR_CHARACTER_LIMIT} characters`}
              </HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
      </FormGroup>
      <FormGroup label="Model description" fieldId="model-description">
        <FormFieldset component={modelDescriptionInput} field="Model Description" />
        <FormHelperText>
          <HelperText>
            <HelperTextItem>Enter a brief summary of the model&apos;s key details.</HelperTextItem>
          </HelperText>
        </FormHelperText>
      </FormGroup>
    </FormSection>
  );
};

export default RegisterModelDetailsFormSection;
