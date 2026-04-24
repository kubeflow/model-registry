import { FormGroup } from '@patternfly/react-core';
import React from 'react';
import { SimpleSelect } from 'mod-arch-shared';
import { SimpleSelectOption } from 'mod-arch-shared/dist/components/SimpleSelect';
import { ModelRegistryCustomProperties } from '~/app/types';
import { ModelType } from '~/concepts/modelCatalog/const';
import { formatModelTypeDisplay } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import {
  buildCustomPropertiesWithModelType,
  getModelTypeStoredValueFromCustomProperties,
} from './registerModelTypeUtils';

const MODEL_TYPE_SELECT_OPTIONS: SimpleSelectOption[] = [
  {
    key: ModelType.GENERATIVE,
    label: formatModelTypeDisplay(ModelType.GENERATIVE),
  },
  {
    key: ModelType.PREDICTIVE,
    label: formatModelTypeDisplay(ModelType.PREDICTIVE),
  },
];

type RegisterModelTypeFieldProps = {
  modelCustomProperties: ModelRegistryCustomProperties | undefined;
  onModelCustomPropertiesChange: (next: ModelRegistryCustomProperties) => void;
  isRequired?: boolean;
};

const RegisterModelTypeField: React.FC<RegisterModelTypeFieldProps> = ({
  modelCustomProperties,
  onModelCustomPropertiesChange,
  isRequired,
}) => {
  const stored = getModelTypeStoredValueFromCustomProperties(modelCustomProperties);

  const handleChange = (key: string) => {
    if (key === ModelType.GENERATIVE || key === ModelType.PREDICTIVE) {
      onModelCustomPropertiesChange(buildCustomPropertiesWithModelType(modelCustomProperties, key));
    }
  };

  return (
    <FormGroup label="Model type" isRequired={isRequired} fieldId="register-model-type">
      <FormFieldset
        field="Model type"
        component={
          <SimpleSelect
            options={MODEL_TYPE_SELECT_OPTIONS}
            value={stored ?? undefined}
            onChange={handleChange}
            placeholder="Select model type"
            isFullWidth
            dataTestId="register-model-type-select"
            previewDescription={false}
            popperProps={{ direction: 'down' }}
            toggleProps={{ id: 'register-model-type-toggle' }}
          />
        }
      />
    </FormGroup>
  );
};

export default RegisterModelTypeField;
