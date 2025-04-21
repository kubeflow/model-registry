import React from 'react';
import { FormGroup, TextInput } from '@patternfly/react-core';
import { TypeaheadSelect } from 'mod-arch-shared';
import { TypeaheadSelectOption } from 'mod-arch-shared/dist/components/TypeaheadSelect';
import { RegisteredModel } from '~/app/types';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

type RegisteredModelSelectorProps = {
  registeredModels: RegisteredModel[];
  registeredModelId: string;
  setRegisteredModelId: (id: string) => void;
  isDisabled: boolean;
};

const RegisteredModelSelector: React.FC<RegisteredModelSelectorProps> = ({
  registeredModels,
  registeredModelId,
  setRegisteredModelId,
  isDisabled,
}) => {
  const options: TypeaheadSelectOption[] = React.useMemo(
    () =>
      registeredModels.map(({ name, id }) => ({
        content: name,
        value: id,
        isSelected: id === registeredModelId,
      })),
    [registeredModels, registeredModelId],
  );

  const modelNameInput = (
    <TextInput
      isDisabled
      isRequired
      type="text"
      id="model-name"
      name="registered-model-prefilled"
      value={options.find(({ value }) => value === registeredModelId)?.content}
    />
  );

  if (isDisabled && registeredModelId) {
    /*
      If we're registering a new version for an existing model, we prefill the model and don't allow it to change.
      TODO: We should just be using the `isDisabled` prop of TypeaheadSelect instead of a separate disabled text field,
        but TypeaheadSelect doesn't currently have a way to prefill the selected item / lift the selection state.
        See related PatternFly issue https://github.com/patternfly/patternfly-react/issues/10842
    */
    return (
      <FormGroup label="Model name" className="form-group-disabled" isRequired fieldId="model-name">
        <FormFieldset component={modelNameInput} field="Model Name" />
      </FormGroup>
    );
  }

  return (
    <TypeaheadSelect
      id="model-name"
      onClearSelection={() => setRegisteredModelId('')}
      selectOptions={options}
      isScrollable
      placeholder="Select a registered model"
      noOptionsFoundMessage={(filter) => `No results found for "${filter}"`}
      onSelect={(_event, selection) => {
        setRegisteredModelId(String(selection));
      }}
    />
  );
};

export default RegisteredModelSelector;
