import React from 'react';
import { TextInput } from '@patternfly/react-core';
import { TypeaheadSelect, TypeaheadSelectOption } from '@patternfly/react-templates';
import { RegisteredModel } from '~/app/types';

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

  if (isDisabled && registeredModelId) {
    /*
      If we're registering a new version for an existing model, we prefill the model and don't allow it to change.
      TODO: We should just be using the `isDisabled` prop of TypeaheadSelect instead of a separate disabled text field,
        but TypeaheadSelect doesn't currently have a way to prefill the selected item / lift the selection state.
        See related PatternFly issue https://github.com/patternfly/patternfly-react/issues/10842
    */
    return (
      <TextInput
        isDisabled
        isRequired
        type="text"
        id="model-name"
        name="registered-model-prefilled"
        value={options.find(({ value }) => value === registeredModelId)?.content}
      />
    );
  }

  return (
    <TypeaheadSelect
      id="model-name"
      initialOptions={options}
      placeholder="Select a registered model"
      noOptionsFoundMessage={(filter) => `No results found for "${filter}"`}
      onSelect={(_event, selection) => {
        setRegisteredModelId(String(selection));
      }}
    />
  );
};

export default RegisteredModelSelector;
